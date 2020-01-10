package sensor

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

/**
 * which attached to the device ips
 */
func (dl *LocalDeviceDetail) GetLocalSensorList(attachIP string) []LocalSensorInformation {
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
	// Interval int    // 最小间隔时间(如果测量没有被阻塞的话将会接近的最小间隔时间, 阻塞的原因可能归结于透传设备无法响应请求指令)
}

type TaskSensorBody struct {
	TaskSensorKey TaskSensorKey // 任务唯一id
	Type          byte          // 指令类型
	RequestData   []byte        // 生成的指令数据
	SensorID      string        // 传感器ID

	customFunction func(body TaskSensorBody, wg *sync.WaitGroup)
}

var tw *TimeWheel

const taskSecond int64 = 1000000000

// 指令类型Type
const DissolvedOxygenAndTemperature byte = 0x01 // 溶氧量和温度
const D2 byte = 0x02                            // 未定义的类型
const D3 byte = 0x04                            // ..
const D4 byte = 0x08                            // ..
const D5 byte = 0x10                            // ..
const D6 byte = 0x20                            // ..
const D7 byte = 0x40                            // ..
const D8 byte = 0x80                            // 未定义的类型

// enum
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
func SensorDefaultHandler(body TaskSensorBody, wg *sync.WaitGroup) {

	switch body.Type {
	case DissolvedOxygenAndTemperature:
		// 得到透传conn
		b, _ := GetDeviceSession(body.TaskSensorKey.Attach)
		// 合成地址
		body.CreateMeasureRequest()
		fmt.Printf("[INFO] 测量请求 -> 传感器设备ID %s | 设备地址 %d | 任务类型 %d | 请求数据 %b |\n", body.SensorID, body.TaskSensorKey.Addr, body.Type, body.RequestData)
		// 向传感器发送对应测量请求
		p, err := b.MeasureRequest(body.RequestData, []string{"Oxygen", "Temp"})
		if err != nil {
			fmt.Println("[FAIL] 请求失败")
			// TODO 当DTU没有接受到传感器响应时, 决定是否要结束
			return
		}
		p.SensorID = body.SensorID
		send, err := json.Marshal(p)
		client, _ := GetMQTTInstance()
		client.Publish("sensor/oxygen/measure", 1, false, send)
		break
	case D2:
		// TODO
		fmt.Println("d2")
		break
	case D3:
		// TODO
		fmt.Println("d3")
		break
	case D4:
		// TODO
		fmt.Println("d4")
		break
	case D5:
		// TODO
		fmt.Println("d5")
		break
	case D6:
		// TODO
		fmt.Println("d6")
		break
	case D7:
		// TODO
		fmt.Println("d7")
		break
	case D8:
		// TODO
		fmt.Println("d8")
		break
	default:
		fmt.Println("default")
	}
	wg.Done()
}

/**
 * 使用自定义Handler以代替默认处理过程
 */
func (ls *LocalSensorInformation) AddTaskHandler(callback func(body TaskSensorBody, wg *sync.WaitGroup)) {
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
 * @param queueChannel 单DTU内任务的阻塞队列
 * -1 -> 无限次
 * >1 -> 有限次
 * @return error 错误的添加会触发
 */
func (ls *LocalSensorInformation) CreateTask(times int, queueChannel chan TaskSensorBody) error {
	key := TaskSensorKey{ls.Addr, ls.Attach, ls.Type}

	body := TaskSensorBody{}
	body.TaskSensorKey = key
	body.Type = ls.Type
	body.RequestData = nil
	body.SensorID = ls.SensorID
	if ls.TaskHandler == nil {
		body.customFunction = nil
	} else {
		body.customFunction = ls.TaskHandler
	}
	data := TaskData{"Data": body, "Channel": queueChannel}
	return tw.AddTask(time.Duration(ls.Interval*taskSecond), times, key, data, TaskSensorPush)
}

/**
 * 单DTU任务阻塞队列的压入回调
 * @param queueChannel 单DTU内任务的阻塞队列
 */
func TaskSensorPush(data TaskData) {
	body := data["Data"].(TaskSensorBody)
	queueChannel := data["Channel"].(chan TaskSensorBody)
	queueChannel <- body
}

/**
 * 单DTU任务调度Routine
 * @param queueChannel 单DTU内任务的阻塞队列
 */
func TaskSensorPop(queueChannel chan TaskSensorBody) {
	var wg sync.WaitGroup
	for {
		v, ok := <-queueChannel
		if !ok {
			break
		} else {
			wg.Add(1)
			SensorDefaultHandler(v, &wg)
			wg.Wait()
		}
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
func (ls *LocalSensorInformation) UpdateTask(interval time.Duration, queueChannel chan TaskSensorBody) error {
	key := TaskSensorKey{ls.Addr, ls.Attach, ls.Type}
	body := TaskSensorBody{}
	body.TaskSensorKey = key
	body.Type = ls.Type
	body.RequestData = nil
	body.SensorID = ls.SensorID
	if ls.TaskHandler == nil {
		body.customFunction = nil
	} else {
		body.customFunction = ls.TaskHandler
	}
	data := TaskData{"Data": body, "Channel": queueChannel}
	return tw.UpdateTask(key, interval, data)
}

/**
 * 初始化TimeWheel
 */
func TimeWheelInit() *TimeWheel {
	tw = New(time.Second, 180)
	tw.Start()
	return tw
}

/*
 * 获得TimeWheel指针
 */
func GetTimeWheel() *TimeWheel {
	return tw
}

func TaskSetup(ip string) {
	ch := make(chan TaskSensorBody, 10)
	go TaskSensorPop(ch)
	// 此处得到attach到该dtu的至少0个, 至多3个传感器的参数
	// TODO: 重构1 初始化传感器状态
	for _, v := range GetLocalDevices().GetLocalSensorList(ip) {
		if err := v.CreateTask(-1, ch); err != nil {
			continue
		}
		fmt.Printf("[INFO] ID:%s 进入队列\n", v.SensorID)
	}
}
