package sensor

import (
	"gopkg.in/mgo.v2/bson"
)

type DeviceOperation interface {
	// 出厂设置
	RestoreFactory()

	// 传感器地址相关
	GetSensorAddr() []byte
	SetSensorAddr()

	// 校准相关
	// Zero and Gradient
	GetCorrectValue()
	SetCorrectValue()

	// 数据获取相关
	GetMeasuredValue()

	SetDefault()

	// 请求结构构件
	RequestBuilder()
}

type Sensor struct {
	// the sensor address range for 0 - 255, using to verify the sensor and then bind with VID.
	SensorAddr uint8
	// VID is a virtual id assign the sensor address and store in database.
	VID string
}

func (i Sensor) RestoreFactory() {
	panic("implement me")
}

func (i Sensor) GetSensorAddr() []byte {
	return []byte{0x00, 0x01}
}

func (i Sensor) SetSensorAddr() {
	panic("implement me")
}

func (i Sensor) GetCorrectValue() {
	panic("implement me")
}

func (i Sensor) SetCorrectValue() {
	panic("implement me")
}

func (i Sensor) GetMeasuredValue() {
	panic("implement me")
}

func (i Sensor) SetDefault() {
	panic("implement me")
}

// Measure(read)
func Measure(Data []byte, callback func(meta interface{}, data []byte)) {

}

// Config(write)
func Config(Data []byte, callback func(meta interface{}, data []byte)) {

}

func MeasureRequest(id bson.ObjectId, funcCode []byte, callback func(meta interface{}, data []byte)) {

}

func ConfigRequest(id bson.ObjectId, funcCode []byte, callback func(meta interface{}, data []byte)) {

}

func AddDevice(id bson.ObjectId) {

}

func RemoteDevice(id bson.ObjectId) {

}
