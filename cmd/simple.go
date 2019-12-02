/**
 * 一个简单的启动示例
 * @data 2019/11/23
 */
package main

import (
	"fmt"
	"sensor"
	"time"
)

func main() {
	list := sensor.LoadConfig("cnf/conf.json")
	fmt.Println(list)
	sensor.TimeWheelInit()
	for _, v := range list.GetLocalSensorList("192.168.5.55") {
		if err := v.CreateTask(time.Second*5, -1); err != nil {
			continue
		}
		fmt.Println("[INFO] 添加一个新的测量任务")
	}

	time.Sleep(time.Minute * 3)
}
