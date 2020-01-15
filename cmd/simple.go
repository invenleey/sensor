/**
 * 一个简单的启动示例
 * @data 2019/11/23
 * @update 2019/12/6
 */
package main

import (
	"sensor"
	"time"
)

func main() {
	// 订阅示例
	// 上级: GO -> MQTT -> GO

	// 测量
	sensor.MQTTMapping("sensor/oxygen/measure", sensor.FncTest)

	// CONFIG更新(重启生效)
	sensor.MQTTMapping("sensor/setting/all", sensor.SettingConfigHandler)

	// Status&&Exception动态更新
	sensor.MQTTMapping("sensor/action/clear", sensor.ClearExceptionHandler)

	// 重启
	sensor.MQTTMapping("sensor/action/restart", sensor.RestartHandler)

	time.AfterFunc(time.Second*20, func() {
		client, _ := sensor.GetMQTTInstance()
		client.Publish("sensor/action/restart", 1, false, "1111")
	})

	// 服务开启示例
	// 下级: GO -> DTU -> Sensor
	go sensor.RunDeviceTCP()

	// token.Wait()

	time.Sleep(time.Hour)

}
