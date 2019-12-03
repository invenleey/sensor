package sensor

import (
	"encoding/json"
	"fmt"
	"time"
)

/**
 * which attached to the device ips
 */
func (dl *LocalDeviceList) GetLocalSensorList(attachIP string) []LocalSensorInformation {
	var ret []LocalSensorInformation
	for _, v := range dl.LocalSensorInformation {
		if v.Attach == attachIP {
			ret = append(ret, v)
		}
	}
	return ret
}

func (ls *LocalSensorInformation) StartSensorMeasureTask() {

}

type TaskSensorKey struct {
	Addr   byte   // 设备地址
	Attach string // 附着设备Gateway
	Type   byte   // 指令类型
	// Interval int    // 最大间隔时间
}

type TaskSensorBody struct {
	TaskSensorKey TaskSensorKey // 任务唯一id
	Type          byte          // 指令类型
	RequestData   []byte        // 生成的指令数据
}

var tw *TimeWheel

const taskSecond int64 = 1000000000

// 指令类型Type
const DissolvedOxygenAndTemperature byte = 0x01
const D2 byte = 0x02
const D3 byte = 0x04
const D4 byte = 0x08
const D5 byte = 0x10
const D6 byte = 0x20
const D7 byte = 0x40
const D8 byte = 0x80

// ...
// ...
// ...

/**
 * 测量请求体创建
 */
func (ts *TaskSensorBody) CreateMeasureRequest() {
	var sr []byte
	// 设备ADDR
	sr = append(sr, ts.TaskSensorKey.Addr)
	// 指令功能码
	sr = append(sr, InfoMK["ReadFunc"]...)
	// 寄存器地址和数量
	sr = append(sr, InfoMK["RMeasure"]...)
	// CRC_ModBus
	sr = append(sr, CreateCRC(sr)...)
	ts.RequestData = sr
}

/**
 * 默认处理过程
 * 当LocalSensorInformation没有设置handler时, 所调用的默认处理过程
 * DefaultHandler中规定了几种默认的处理方式
 */
func sensorDefaultHandler(data TaskData) {
	body := data["Data"].(TaskSensorBody)
	fmt.Printf("[INFO] 设备地址 %d 任务类型 %d 施工中\n", body.TaskSensorKey.Addr, body.Type)

	switch body.Type {
	case DissolvedOxygenAndTemperature:
		fmt.Println("溶氧量和温度查询过程")
		// 得到透传conn
		b, _ := GetDeviceSession(body.TaskSensorKey.Attach)
		// 合成地址
		body.CreateMeasureRequest()
		fmt.Println(body.RequestData)
		// 向传感器发送对应测量请求
		p, err := b.MeasureRequest(body.RequestData, []string{"测量值", "温度"})
		// fmt.Println(body.RequestData)
		//p, err = b.MeasureRequest(body.RequestData, []string{"aas", "bbs"})
		if err != nil {
			fmt.Println("[FAIL] 请求失败")
			fmt.Println(err)
		}
		send, err := json.Marshal(p)
		client, _ := GetMQTTInstance()
		client.Publish("sensor/oxygen/measure", 1, false, send)
		break
	case D2:
		fmt.Println("d2")
		break
	case D3:
		fmt.Println("d3")
		break
	case D4:
		fmt.Println("d4")
		break
	case D5:
		fmt.Println("d5")
		break
	case D6:
		fmt.Println("d6")
		break
	case D7:
		fmt.Println("d7")
		break
	case D8:
		fmt.Println("d8")
		break
	default:
		fmt.Println("default")
	}
}

/**
 * 使用自定义Handler以代替默认处理过程
 */
func (ls *LocalSensorInformation) AddTaskHandler(callback Job) {
	ls.TaskHandler = callback
}

/**
 * 移除自定义Handler
 */
func (ls *LocalSensorInformation) RemoveTaskHandler() bool {
	if ls.TaskHandler == nil {
		return false
	}
	ls.TaskHandler = nil
	return true
}

/**
 * 启动传感器任务
 * @param ls.interval 测量间隔时间
 * @param times 指定任务次数
 * -1 -> 无限次
 * >1 -> 有限次
 * @return error 错误的添加会触发
 */
func (ls *LocalSensorInformation) CreateTask(times int) error {
	key := TaskSensorKey{ls.Addr, ls.Attach, ls.Type}
	body := TaskSensorBody{key, ls.Type, nil}
	data := TaskData{"Data": body}
	if ls.TaskHandler == nil {
		return tw.AddTask(time.Duration(ls.Interval*taskSecond), times, key, data, sensorDefaultHandler)
	} else {
		return tw.AddTask(time.Duration(ls.Interval*taskSecond), times, key, data, ls.TaskHandler)
	}
}

/**
 * 移除传感器任务
 * @return error 错误的移除同样会触发
 */
func (ls *LocalSensorInformation) RemoveTask() error {
	key := TaskSensorKey{ls.Addr, ls.Attach, ls.Type}
	return tw.RemoveTask(key)
}

/**
 * 更新传感器任务
 * @param interval 测量间隔时间
 * @return error 需要注意多次提交相同key任务时会触发
 *
 */
func (ls *LocalSensorInformation) UpdateTask(interval time.Duration) error {
	key := TaskSensorKey{ls.Addr, ls.Attach, ls.Type}
	body := TaskSensorBody{key, ls.Type, nil}
	data := TaskData{"Data": body}
	return tw.UpdateTask(key, interval, data)
}

/**
 * 初始化
 */
func TimeWheelInit() *TimeWheel {
	tw = New(time.Second, 180)
	tw.Start()
	return tw
}

/*
 * 获得
 */
func GetTimeWheel() *TimeWheel {
	return tw
}
