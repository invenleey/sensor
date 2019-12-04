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
	ips := sensor.ShowNodeIPs()
	for _, ip := range ips {
		for _, v := range list.GetLocalSensorList(ip) {

			if err := v.CreateTask(-1); err != nil {
				continue
			}
			fmt.Printf("[INFO] ID:%s 进入队列", v.SensorID)
		}
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
