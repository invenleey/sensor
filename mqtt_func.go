package sensor

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"sensor/count"
)

/**
 * 测试
 */
var FncTest = func(client mqtt.Client, message mqtt.Message) {
	fmt.Printf("主题: %s\n", message.Topic())
	fmt.Printf("信息: %s\n", message.Payload())
}

/**
 * 以传感器为单位的基本构成如下功能
 * 1. 修改测量间隔时间
 * 2. 修改关联IP
 * 3. ...暂定这两个
 *
 * 即含
 * 1. 对象 sensorID
 * 2. 操作 operation (type)
 * 3. 数据 data
 */
type SensorAction struct {
	SensorID  string      `json:"sensorID"`
	Operation string      `json:"operation"`
	Data      interface{} `json:"data"`
}

/*
 * 清除错误次数
 * @Topic sensor/action/clear
 *
 */
func ClearExceptionHandler(client mqtt.Client, message mqtt.Message) {
	sa, _ := RequestMap(message)
	switch sa.Operation {
	case count.CLEAR_ALL_EXCEPTION:
		count.ClsAll()
		break
	case count.CLEAR_ONE_EXCEPTION:
		count.ClsErrorCount(sa.SensorID)
		break
	default:
	}
}

/**
 * 设置/更改传感器关联Attach
 * @Topic sensor/action/setAttachIP
 *
 */
func ChangeAttachIPHandler(client mqtt.Client, message mqtt.Message) {
	sa, _ := RequestMap(message)
	if IsIp(sa.Data.(string)) {

	}
}

/*
 * 更改当前节点上的CONFIG文件
 *
 */
//func ()  {
//
//}

/**
 * action字段映射
 */
func RequestMap(message mqtt.Message) (*SensorAction, error) {
	var sensorAction SensorAction
	err := json.Unmarshal(message.Payload(), &sensorAction)
	return &sensorAction, err
}
