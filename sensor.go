package sensor

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"
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
var InfoMK = map[string][]byte{
	"ReadFunc":  {0x03},
	"WriteFunc": {0x06},

	"RMeasure": {0x00, 0x00, 0x00, 0x04},
	"WOxygen":  {0x00, 0x04, 0x00, 0x01},
	"WZero":    {0x10, 0x00, 0x00, 0x01},
	"WTilt":    {0x10, 0x04, 0x00, 0x01},

	"RZero": {010, 0x06, 0x00, 0x01},
	"RTilt": {0x10, 0x08, 0x00, 0x01},

	"RAddr":    {0x20, 0x02, 0x00, 0x01},
	"WFactory": {0x20, 0x20, 0x00, 0x01},
}

//var request []byte

//// Measure(read)
//func Measure(addr []byte, callback func(meta interface{}, data []byte)) {
//
//}

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

/**
 * the measure values struct
 */
type ReadResult struct {
	SensorID string `json:"sensorID"`
	// unique DeviceID in node server
	DeviceAddr byte `json:"deviceAddr"`
	// Function Code which had been operate
	FuncCode byte `json:"funcCode"`
	// Information Count
	InfoCount int `json:"infoCount"`

	// read
	Items []MeasureItem `json:"items,omitempty"`

	// order
	WriteData []byte `json:"writeData,omitempty"`
	WriteReg  []byte `json:"writeReg,omitempty"`

	// node server ip
	NodeIP string `json:"nodeIP"`

	// create time
	Created time.Time `json:"created"`

	// error tag
	Status int `json:"status"`
}

type MeasureItem struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

/**
 * @return a measureResult data struct
 */
func (ds *DeviceSession) GetResultInstance(meta DeviceMeta) (ReadResult, error) {
	var ins ReadResult

	ins.Created = time.Now()
	ins.DeviceAddr = meta.Addr
	ins.FuncCode = meta.FuncCode
	ins.NodeIP = ds.conn.RemoteAddr().String()
	// check order whether is wrong
	if meta.FuncCode > 0x80 {
		ins.Status = 1
		ins.FuncCode -= 0x80
		return ins, errors.New("unknown order")
	}
	return ins, nil
}

/**
 * decode data
 * @param data the measure data which will be decoded
 * @param df is the deviceAddr and funcCode
 */
func (mr *ReadResult) DecodeStandardFourByte2Float(data []byte, itemsName []string) error {
	v, err := FourByteToFloat(data)
	if err != nil {
		return err
	}
	if len(v) != len(itemsName) {
		return errors.New("error itemsName count")
	}

	// inject content
	mr.InfoCount = len(v)
	for i, k := range v {
		var item MeasureItem
		item.Name = itemsName[i]
		item.Value = k
		mr.Items = append(mr.Items, item)
	}
	return nil
}

func (mr *ReadResult) DecodeSlope(data []byte, itemName string) error {
	v, err := TwoByteToFloatX1000(data)
	if err != nil {
		return errors.New("get error type")
	}
	var item MeasureItem
	item.Name = itemName
	item.Value = v
	mr.Items = append(mr.Items, item)
	return nil
}

// ====================================Write Process======================================== //
//type WriteResult struct {
//	// unique DeviceID in node server
//	DeviceAddr byte
//	// Function Code which had been operate
//	FuncCode byte
//
//
//	// node server ip
//	NodeIP string
//
//	// error tag
//	// 0 succeed
//	// 1 failed
//	status int
//}

func (mr *ReadResult) DecodeOrder(data []byte) error {
	mr.WriteData = data[2:]
	mr.WriteReg = data[:2]
	return nil
}
