package sensor

import (
	"encoding/json"
	"testing"
	"time"
)

func TestSettingConfigHandler(t *testing.T) {
	MQTTMapping("sensor/setting/all", SettingConfigHandler)

	i := SensorAction{
		SensorID:  "abcdefg",
		Operation: "del",
	}

	v := GetLocalDevicesInstance()
	i.Data, _ = json.Marshal(v)

	j, _ := json.Marshal(i)
	client, _ := GetMQTTInstance()
	token := client.Publish("sensor/setting/all", 1, false, j)
	token.Wait()

	time.Sleep(time.Minute)

}

func TestRestartHandler(t *testing.T) {
	client, _ := GetMQTTInstance()
	client.Publish("sensor/action/bbbb", 1, false, 1)


}
