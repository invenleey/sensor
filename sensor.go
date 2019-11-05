package sensor

import (
	"errors"
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

//// Measure(read)
//func Measure(Data []byte, callback func(meta interface{}, data []byte)) {
//
//}
//
//// Config(write)
//func Config(Data []byte, callback func(meta interface{}, data []byte)) {
//
//}

//func MeasureRequest(id bson.ObjectId, funcCode []byte, callback func(meta interface{}, data []byte)) {
//
//}
//
//func ConfigRequest(id bson.ObjectId, funcCode []byte, callback func(meta interface{}, data []byte)) {
//
//}

// Freedom map
var InfoMK map[string][]byte

func InitInfoMK() {
	// Function Code Type Map
	// Read 0x03
	// Write 0x06
	InfoMK["ReadFunc"] = []byte{0x03}
	InfoMK["WriteFunc"] = []byte{0x06}

	InfoMK["RRegMeasure"] = []byte{0x00, 0x00, 0x00, 0x04}
	InfoMK["WRegOxygen"] = []byte{0x10, 0x04}
	InfoMK["WRegZero"] = []byte{0x10, 0x00}
	InfoMK["WRegTilt"] = []byte{0x10, 0x04}

	InfoMK["RRegZero"] = []byte{010, 0x06}
	InfoMK["RRegTilt"] = []byte{0x10, 0x08}

	InfoMK["ARegAddr"] = []byte{0x20, 0x02}
	InfoMK["WRegFactory"] = []byte{0x20, 0x20}

}

// Function Code Type Map
// Read 0x03
// Write 0x06
//var ReadFunc = []byte{0x03}
//var WriteFunc = []byte{0x06}
//
//// Register Address include reg count
//var RRegMeasure = []byte{0x00, 0x00, 0x00, 0x04}
//var WRegOxygen = []byte{0x10, 0x04}
//var WRegZero = []byte{010, 0x00}
//var WRegTilt = []byte{0x10, 0x04}
//
//var RRegZero = []byte{010, 0x06}
//var RRegTilt = []byte{0x10, 0x08}
//
//var ARegAddr = []byte{0x20, 0x02}
//var WRegFactory = []byte{0x20, 0x20}
//
//var request []byte

// Measure(read)
func Measure(addr []byte, callback func(meta interface{}, data []byte)) {

}

// Config(write)
func Config(Data []byte, callback func(meta interface{}, data []byte)) {

}

//func MeasureRequest(addr string, callback func(meta interface{}, data []byte)) {
//	da := ComposeBody([]byte{0x06}, ReadFunc, RRegMeasure)
//	fmt.Println("test: ", da)
//	// SendWord(addr, da, callback)
//}

func ConfigRequest(id bson.ObjectId, funcCode []byte, callback func(meta interface{}, data []byte)) {

}

func AddDevice(id bson.ObjectId) {

}

func RemoteDevice(id bson.ObjectId) {

}

/**
 * the measure values struct
 */
type MeasureResult struct {
	// unique DeviceID in node server
	DeviceAddr byte
	// Function Code which had been operate
	FuncCode byte
	// Information Count
	InfoCount int

	Items []MeasureItem

	// node server ip
	NodeIP string
}

type MeasureItem struct {
	Name  string
	Value float64
}

/**
 * @return a measureResult data struct
 */
func (ds *DeviceSession) GetMeasureResultInstance() MeasureResult {
	var ins MeasureResult
	ins.NodeIP = ds.conn.RemoteAddr().String()
	return ins
}

/**
 * decode data
 * @param data the measure data which will be decoded
 * @param df is the deviceAddr and funcCode
 */
func (mr *MeasureResult) DecodeMeasureByte(meta DeviceMeta, data []byte, funcCode byte, itemsName []string) error {
	v, err := FourByteToFloat(data)
	if err != nil {
		return err
	}
	if len(v) != len(itemsName) {
		return errors.New("error itemsName count")
	}
	if meta.FuncCode == funcCode+0x80 {
		if data[0] == 0x01 {
			mr.DeviceAddr = meta.Addr
			mr.FuncCode = meta.FuncCode
			return errors.New("error v-order")
		} else {
			return errors.New("error v-data")
		}
	}
	// inject content
	mr.DeviceAddr = meta.Addr
	mr.FuncCode = meta.FuncCode
	mr.InfoCount = len(v)
	for i, k := range v {
		var item MeasureItem
		item.Name = itemsName[i]
		item.Value = k
		mr.Items = append(mr.Items, item)
	}
	return nil
}
