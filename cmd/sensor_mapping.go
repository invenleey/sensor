/**
 * 一个简单的启动示例
 * @data 2019/11/23
 * @update 2019/12/6
 */
package main

import (
	"sensor"
)

func SensorMapping() {
	// 订阅示例: 下位 -> MQTT -> 上位

	// 测量
	sensor.MQTTMapping("sensor/oxygen/measure", sensor.FncTest)

	// CONFIG更新(重启生效)
	sensor.MQTTMapping("sensor/setting/all", sensor.SettingConfigHandler)

	// Status&&Exception动态更新
	sensor.MQTTMapping("sensor/action/clear", sensor.ClearExceptionHandler)

	// 重启
	sensor.MQTTMapping("sensor/action/restart", sensor.RestartHandler)

	// 状态开关
	sensor.MQTTMapping("sensor/action/switch", sensor.SwitchSensorHandler)

}
