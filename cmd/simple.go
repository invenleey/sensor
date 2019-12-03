/**
 * 一个简单的启动示例
 * @data 2019/11/23
 */
package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"sensor"
	"time"
)

func main() {
	// TCP server
	go sensor.RunDeviceTCP()

	// 暂停等待测试
	time.Sleep(time.Second * 10)

	// 读取配置
	list := sensor.LoadConfig("cnf/conf.json")
	fmt.Println(list)
	// 初始化
	sensor.TimeWheelInit()
	for _, v := range list.GetLocalSensorList("172.20.10.4") {

		if err := v.CreateTask(-1); err != nil {
			continue
		}
		fmt.Println("[INFO] DEMO 添加一个新的测量任务")
	}

	ins, _ := sensor.GetMQTTInstance()
	if token := ins.Subscribe("sensor/oxygen/measure", 1, f); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
	time.Sleep(time.Minute * 3)
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}