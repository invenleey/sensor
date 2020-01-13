package sensor

import (
	"encoding/json"
	"fmt"
)

//=====================LOG======================
//
//                   日志上报
//
//=====================END======================

const (
	LOG_INFO    = "[INFO] "
	LOG_WARN    = "[WARN] "
	LOG_FAIL    = "[FAIL] "
	MQ_LOG_INFO = iota
	MQ_LOG_WARN
	MQ_LOG_FAIL
)

type MQLog struct {
	SensorID   string // 传感器ID(可无)
	LogMessage string // 日志内容
	LogLevel   int    // 日志登记

}

// 限制等级
var messageLevel = MQ_LOG_FAIL

func SetMQLogLevel(level int) {
	messageLevel = level
}

func PushMQLog(logLevel int, msg string, sensorID ...string) {

	var mql MQLog
	if len(sensorID) > 0 {
		mql.SensorID = sensorID[0]
	}
	mql.LogLevel = logLevel
	mql.LogMessage = msg
	mql.Printf()
}

/**
 *
 */
func (ml *MQLog) Printf() {
	var msg string
	switch ml.LogLevel {
	case MQ_LOG_INFO:
		msg = LOG_INFO + ml.LogMessage
		break
	case MQ_LOG_WARN:
		msg = LOG_WARN + ml.LogMessage
		break
	case MQ_LOG_FAIL:
		msg = LOG_FAIL + ml.LogMessage
		break
	default:
	}

	fmt.Println(msg)

	if messageLevel > ml.LogLevel {
		return
	}

	cli, _ := GetMQTTInstance()
	jml, _ := json.Marshal(ml)
	cli.Publish("sensor/log", 1, false, jml)
}
