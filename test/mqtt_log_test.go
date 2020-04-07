package test

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"sensor"
	"testing"
	"time"
)

func TestPushMQLog(t *testing.T) {
	sensor.MQTTMapping("sensor/log", func(c mqtt.Client, message mqtt.Message) {
		fmt.Println(message.Payload())
	})

	sensor.PushMQLog(sensor.MQ_LOG_FAIL, "hello")

	time.Sleep(time.Minute)
}
