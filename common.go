package sensor

import (
	"bytes"
	"encoding/binary"
	"errors"
	"gopkg.in/mgo.v2/bson"
)

/**
 * bigEndian transfer
 */
func ToBigEndian(num uint16) []byte {
	int16buf := new(bytes.Buffer)
	if err := binary.Write(int16buf, binary.BigEndian, num); err != nil {
		panic("error type: num")
	}
	return int16buf.Bytes()
}

/**
 * littleEndian transfer
 */
func ToLittleEndian(num uint16) []byte {
	int16buf := new(bytes.Buffer)
	if err := binary.Write(int16buf, binary.LittleEndian, num); err != nil {
		panic("error type: num")
	}
	return int16buf.Bytes()
}

/**
 * bytes to uint(bigEndian)
 */
func BytesToIntU(b []byte) (uint16, error) {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp uint16
	err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return tmp, err
}

/**
 * generate a default request command
 */
func ComposeBody(DeviceAddr, FuncCode, Data []byte) []byte {
	ret := append(DeviceAddr, FuncCode...)
	ret = append(ret, Data...)
	ret = append(ret, CreateCRC(ret)...)
	return ret
}

/**
 * separate config command parameters
 * @params src is a measure respond data
 * @return DeviceAddr and FuncCode
 * @return RegisterAddr
 * @return ConfigData
 * @return error it will also use crc-16 to validate whether it got a true data, if not returns error(wrong data)
 */
func SplitConfig(src []byte) ([]byte, []byte, []byte, error) {
	if ValidateCRC(src[:6], src[6:8]) {
		return src[:2], src[2:4], src[4:6], nil
	}
	return nil, nil, nil, errors.New("got error data")
}

/**
 * separate measure command parameters
 * @params src is a measure respond data
 * @return DeviceAddr and FuncCode
 * @return ByteCount
 * @return MeasureData
 */
func SplitMeasure(src []byte) ([]byte, []byte, error) {
	ByteCount := src[2] + 3
	if ValidateCRC(src[:ByteCount], src[ByteCount:ByteCount+2]) {
		return src[:2], src[3:ByteCount], nil
	}
	return nil, nil, errors.New("got error data")
}

/**
 * get device information from database
 */
func GetDeviceInfo(id bson.ObjectId) {

}
