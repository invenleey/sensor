/**
 * 一个简单的启动示例
 * @data 2019/11/23
 * @update 2019/12/6
 */
package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"sensor"
	mqtt2 "sensor/mqtt"
)

func main() {
	// 订阅示例
	// 上级: GO -> MQTT -> GO
	mqtt2.MQTTMapping("sensor/oxygen/measure", 1,
		func(client mqtt.Client, msg mqtt.Message) {
			fmt.Printf("主题: %s\n", msg.Topic())
			fmt.Printf("信息: %s\n", msg.Payload())
		})

	// 初始化任务轮盘
	// sensor.TimeWheelInit()

	// 服务开启示例
	// 下级: GO -> DTU -> Sensor
	sensor.RunDeviceTCP()
}
