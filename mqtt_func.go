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
	SensorID  string `json:"sensorID"`
	Operation string `json:"operation"`
	Data      []byte `json:"data"`
}

/*
 * 清除错误次数&&状态
 * @Topic sensor/action/clear
 *
 */
func ClearExceptionHandler(client mqtt.Client, message mqtt.Message) {
	sa, _ := RequestMap(message)
	switch sa.Operation {
	case count.CLEAR_ALL_EXCEPTION:
		count.ClsAll()
		for _, v := range GetLocalDevicesInstance().LocalSensorInformation {
			v.Status = STATUS_NORMAL
		}
		break
	case count.CLEAR_ONE_EXCEPTION:
		count.ClsErrorCount(sa.SensorID)
		ld, _ := GetLocalSensor(sa.SensorID)
		ld.Status = STATUS_NORMAL
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
	//sa, _ := RequestMap(message)
	//if IsIp(sa.Data.(string)) {
	//
	//}
}

/*
 * 更改当前下位机上的CONFIG文件
 * Topic sensor/setting/all
 */
func SettingConfigHandler(client mqtt.Client, message mqtt.Message) {
	sa, _ := RequestMap(message)

	// 反序列化
	var jsonConfig LocalDeviceDetail
	if err := json.Unmarshal(sa.Data, &jsonConfig); err != nil {
		fmt.Println("[FAIL] CONFIG反序列化错误")
		return
	}

	// 保存设置
	if err := jsonConfig.DumpConfig(); err != nil {
		fmt.Println("[FAIL] CONFIG保存时错误")
		return
	}

	jsonConfig.ReplaceLocalDeviceInstance()
	RestartDeviceTCP()
	// fmt.Println("[INFO] CONFIG已更新")
}

/**
 * 重启服务 + 重新加载数据
 */
func RestartHandler(client mqtt.Client, message mqtt.Message) {
	fmt.Println("[INFO] 正在重启TCP")
	ReloadDeviceInstance()
	RestartTCPSystem()
}

/**
 * 动态更新
 */
func DynamicAdd(client mqtt.Client, message mqtt.Message) {

}

/**
 * action字段映射
 */
func RequestMap(message mqtt.Message) (*SensorAction, error) {
	var sensorAction SensorAction
	err := json.Unmarshal(message.Payload(), &sensorAction)
	return &sensorAction, err
}
