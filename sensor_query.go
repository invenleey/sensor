package sensor

import (
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
	Addr   int    // 设备地址
	Attach string // 附着设备Gateway
	Type   string // 指令类型
}

type TaskSensorBody struct {
	TaskSensorKey TaskSensorKey // 任务唯一id
	Type          string        // 指令类型
}

var tw *TimeWheel

func sensorMeasureHandler(data TaskData) {
	body := data["Data"].(TaskSensorBody)
	fmt.Printf("[INFO] 设备地址 %d 任务类型 %s 施工中\n", body.TaskSensorKey.Addr, body.Type)
}

/**
 * 启动传感器任务
 * @param interval 测量间隔时间
 * @param times 任务次数
 * -1 -> 无限次
 * >1 -> 有限次
 * @return error 错误的添加会触发
 */
func (ls *LocalSensorInformation) CreateTask(interval time.Duration, times int) error {
	key := TaskSensorKey{ls.Addr, ls.Attach, ls.Type}
	body := TaskSensorBody{key, ls.Type}
	data := TaskData{"Data": body}
	return tw.AddTask(interval, times, key, data, sensorMeasureHandler)
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
	body := TaskSensorBody{key, ls.Type}
	data := TaskData{"Data": body}
	return tw.UpdateTask(key, interval, data)
}

func TimeWheelInit() {
	tw = New(time.Second, 180)
	tw.Start()
}

func GetTimeWheel() *TimeWheel {
	return tw
}
