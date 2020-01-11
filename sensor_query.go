package sensor

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

/**
 * @return 关联至attachIP上的至少0个传感器
 */
func (dl *LocalDeviceDetail) GetLocalSensorList(attachIP string) []*LocalSensorInformation {
	var ret []*LocalSensorInformation
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
	// TaskSensorKey TaskSensorKey // 任务唯一id
	Type           byte   // 指令类型
	RequestData    []byte // 生成的指令数据
	SensorID       string // 传感器ID
	SensorAddr     byte   // 传感器地址
	SensorAttachIP string // 传感器依附IP

	customFunction func(body TaskSensorBody, wg *sync.WaitGroup)
}

var tw *TimeWheel

const taskSecond int64 = 1000000000

const (
	DissolvedOxygenAndTemperature = iota // 溶氧量
	D2
	D3
	D4
	D5
	D6
	D7
	D8
	OnlineScanner // 在线 -> 8
	// 自定义指令Type 用于一次性的用户设置指令等

)

//// 指令类型Type
//const DissolvedOxygenAndTemperature byte = 0x01 // 溶氧量和温度
//const D2 byte = 0x02                            // 未定义的类型
//const D3 byte = 0x04                            // ..
//const D4 byte = 0x08                            // ..
//const D5 byte = 0x10                            // ..
//const D6 byte = 0x20                            // ..
//const D7 byte = 0x40                            // ..
//const D8 byte = 0x80                            // 未定义的类型

// 自定义指令Type 用于一次性的用户设置指令等

// enum
// ...
// ...

/**
 * 测量请求体创建
 */
func (ts *TaskSensorBody) CreateMeasureRequest() {
	var sr []byte
	// 设备ADDR
	sr = append(sr, ts.SensorAddr)
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
func DefaultSensorHandler(body TaskSensorBody, wg *sync.WaitGroup) {

	switch body.Type {
	case DissolvedOxygenAndTemperature:
		// 得到透传conn
		b, _ := GetDeviceSession(body.SensorAttachIP)
		// 合成地址
		body.CreateMeasureRequest()
		fmt.Printf("[INFO] 测量请求 -> 传感器设备ID %s | 设备地址 %d | 任务类型 %d | 请求数据 %b |\n", body.SensorID, body.SensorAddr, body.Type, body.RequestData)
		// 向传感器发送对应测量请求
		p, err := b.MeasureRequest(body.RequestData, []string{"Oxygen", "Temp"})
		if err != nil {
			fmt.Println("[FAIL] 请求失败")
			// TODO 超时处理
			v, _ := GetLocalSensor(body.SensorID)
			// 超时标记
			v.Status = STATUS_DETACH
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
	// key由传感器地址addr + 依附下位机attachIP + 传感器类型type构成
	body := TaskSensorBody{}
	body.SensorAddr = ls.Addr
	body.Type = ls.Type
	body.RequestData = nil
	body.SensorID = ls.SensorID
	body.SensorAttachIP = ls.Attach

	// 是否存在自定义任务
	// 这里应该不需要这个自定义任务了, 应该改到pop中
	//if ls.TaskHandler == nil {
	//	body.customFunction = nil
	//} else {
	//	body.customFunction = ls.TaskHandler
	//}

	// data由信息体data + 阻塞channel构成
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
 * 特别说明: 当遇到的任务不是定时执行的时候, 比如是用户修改了传感器的某一项参数时, 需要提前得知queueChannel的地址
 * @param queueChannel 单DTU内任务的阻塞队列
 */
func TaskSensorPop(queueChannel chan TaskSensorBody) {
	var wg sync.WaitGroup
	for {
		v, ok := <-queueChannel
		if !ok {
			break
		} else {
			// 确保数据有序进行
			wg.Add(1)
			DefaultSensorHandler(v, &wg)
			// 等待上一个任务完成
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

func TaskSetup(attachIP string) {
	ch := make(chan TaskSensorBody, 10)
	go TaskSensorPop(ch)
	// 此处得到attach到该dtu的至少0个, 至多3个传感器的参数
	// TODO: 重构1 初始化传感器状态

	// 为attach的每一个传感器设置定时任务
	for _, v := range GetLocalDevicesInstance().GetLocalSensorList(attachIP) {
		if err := v.CreateTask(-1, ch); err != nil {
			continue
		}
		fmt.Printf("[INFO] ID:%s 进入队列\n", v.SensorID)
	}
	//
}

// 扫描attach(下位机)内传感器状态
// 在processor内的for进行首次判断
func ScanSensorStatus(attach string) {
	// 传感器列表
	sensors := GetLocalDevicesInstance().GetLocalSensorList(attach)
	// session
	ds, _ := GetDeviceSession(attach)
	for _, v := range sensors {
		var sr []byte
		// 设备ADDR
		sr = append(sr, v.Addr)
		// 指令功能码
		sr = append(sr, InfoMK["ReadFunc"]...)
		// 寄存器地址和数量
		sr = append(sr, InfoMK["RAddr"]...)
		// CRC_ModBus
		sr = append(sr, CreateCRC(sr)...)
		if _, err := ds.SendToSensor(sr); err != nil {
			// 超时
			v.Status = STATUS_DETACH
		} else {
			// TODO: 最后记得把fmt换成日志log输出
			v.Status = STATUS_NORMAL
			fmt.Println("[INFO] 设备连接成功" + v.SensorID + " FROM " + v.Attach)
		}
	}
}

// 通过sensorID获得LocalSensorInformation
func GetLocalSensor(sensorID string) (*LocalSensorInformation, error) {
	ins := GetLocalDevicesInstance().LocalSensorInformation
	for _, v := range ins {
		if v.SensorID == sensorID {
			return v, nil
		}
	}
	return nil, errors.New("not find sensorID for this device")
}
