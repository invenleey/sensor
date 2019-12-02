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
	// 读取配置
	list := sensor.LoadConfig("cnf/conf.json")
	fmt.Println(list)
	// 初始化
	sensor.TimeWheelInit()
	for _, v := range list.GetLocalSensorList("192.168.5.55") {

		if err := v.CreateTask(-1); err != nil {
			continue
		}
		fmt.Println("[INFO] DEMO 添加一个新的测量任务")
	}

	time.Sleep(time.Minute * 3)
}
