package test

import (
	"encoding/json"
	"sensor"
	"testing"
	"time"
)

func TestSettingConfigHandler(t *testing.T) {
	sensor.MQTTMapping("sensor/setting/all", sensor.SettingConfigHandler)

	i := sensor.SensorAction{
		SensorID:  "abcdefg",
		Operation: "del",
	}

	v := sensor.GetLocalDevicesInstance()
	i.Data, _ = json.Marshal(v)

	j, _ := json.Marshal(i)
	client, _ := sensor.GetMQTTInstance()
	token := client.Publish("sensor/setting/all", 1, false, j)
	token.Wait()

	time.Sleep(time.Minute)

}

func TestRestartHandler(t *testing.T) {
	client, _ := sensor.GetMQTTInstance()
	client.Publish("sensor/action/bbbb", 1, false, 1)


}
