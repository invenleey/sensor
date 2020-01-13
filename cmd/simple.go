/**
 * 一个简单的启动示例
 * @data 2019/11/23
 * @update 2019/12/6
 */
package main

import (
	"sensor"
)


func main() {
	// 订阅示例
	// 上级: GO -> MQTT -> GO

	sensor.MQTTMapping("sensor/oxygen/measure", sensor.FncTest)

	// 服务开启示例
	// 下级: GO -> DTU -> Sensor
	sensor.RunDeviceTCP()
}
