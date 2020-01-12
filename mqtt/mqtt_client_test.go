package mqtt

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"sensor"
	"testing"
	"time"
)

func TestGetMQTTInstance(t *testing.T) {
	ins, _ := GetMQTTInstance()
	fmt.Println(ins.IsConnectionOpen())
	// drop
	for i := 0; i < 100; i++ {
		text := fmt.Sprintf("this is msg #%d!", i)
		token := ins.Publish("go-mqtt/sample", 2, true, text)
		token.Wait()
	}
	ins.Disconnect(1000)
}

func TestGetInformation(t *testing.T) {
	go sensor.RunDeviceTCP()
	time.Sleep(time.Second * 1)
	var sr []byte
	sr = append(sr, 06)
	sr = append(sr, sensor.InfoMK["ReadFunc"]...)
	sr = append(sr, sensor.InfoMK["RMeasure"]...)
	sr = append(sr, sensor.CreateCRC(sr)...)

	time.Sleep(time.Second * 5)
	b, _ := sensor.GetDeviceSession("192.168.5.94")
	p, _ := b.MeasureRequest(sr, []string{"测量值", "温度"})
	fmt.Println(p)
	fmt.Println(p)
	conn, _ := GetMQTTInstance()

	var cb mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		// fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Println(msg.Payload())
		var m sensor.ReadResult
		_ = json.Unmarshal(msg.Payload(), &m)
		fmt.Println(m)
	}
	conn.Subscribe("iot/temp", 1, cb)
	ba, _ := json.Marshal(p)

	token := conn.Publish("iot/temp", 1, false, ba)
	token.Wait()

	time.Sleep(time.Second * 20)
}
