package sensor

import (
	"dev.atomtree.cn/atom/da"
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

/**
 * measure and config use a same func because they have the third dat which is similar parameter
 */

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

/**
 * this is the iot device struct in database
 */
type IOT struct {
	ID       bson.ObjectId `json:"id" bson:"_id" title:"id"`
	IP       string        `json:"ip" bson:"ip" title:"ip"`
	Pool     bson.ObjectId `json:"pool" bson:"pool" title:"鱼池标识"`
	Number   string        `json:"number" bson:"number" title:"设备编号"`
	Name     string        `json:"name" bson:"name" title:"设备名称"`
	Status   string        `json:"status" bson:"status" title:"在线状态"`
	Detail   string        `json:"detail" bson:"detail" title:"设备详细"`
	Operator string        `json:"operator" bson:"operator" title:"操作员"`
}

/**
 * firstly using ip to identify whether this device is store in database
 * if not, create it on da
 */
func GetDeviceByIP(ip string) (IOT, error) {
	s, db, err := da.ConnectDefault()
	if err != nil {
		return IOT{}, nil
	}
	defer s.Close()

	iotColl := db.C("EventTable")
	iotQuery := bson.M{"ip": ip}

	var iotResult []IOT

	if err = iotColl.Find(iotQuery).Skip(0).Limit(0).Sort().All(&iotResult); err != nil {
		return IOT{}, err
	}
	if len(iotResult) == 0{
		da.save()
	}
	return iotResult, nil
}

