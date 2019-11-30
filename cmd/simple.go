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
	for _, v := range list.GetLocalSensorList("192.168.1.1") {
		if err := v.CreateTask(time.Minute, -1); err != nil {
			continue
		}
		fmt.Println("[INFO] 添加一个新的测量任务")
	}
}
