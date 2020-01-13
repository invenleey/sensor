package sensor

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"testing"
	"time"
)

func TestPushMQLog(t *testing.T) {
	MQTTMapping("sensor/log", func(c mqtt.Client, message mqtt.Message) {
		fmt.Println(message.Payload())
	})

	PushMQLog(MQ_LOG_FAIL, "hello")

	time.Sleep(time.Minute)
}
