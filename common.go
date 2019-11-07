package sensor

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strconv"
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
 * @params ByteCount
 * @return DeviceAddr and FuncCode
 * @return MeasureData
 */
func SplitMeasure(src []byte) (DeviceMeta, []byte, error) {
	base := len(src) - 2
	if ValidateCRC(src[:base], src[base:]) {
		var meta DeviceMeta
		meta.Addr = src[0]
		meta.FuncCode = src[1]
		if src[1] > 0x80 {
			return meta, src[2:base], nil
		} else if src[1] == 0x03 {
			return meta, src[3:base], nil
		} else if src[1] == 0x06 {
			return meta, src[2:base], nil
		}
	}
	return DeviceMeta{}, nil, errors.New("unreachable validate")
}

const HexValue = 256

func ByteToFloat(v []byte) float64 {
	return float64(v[0])*HexValue + float64(v[1])
}

func FourByteToFloat(v []byte) ([]float64, error) {
	count := len(v)
	if count%4 != 0 || count == 0 {
		return nil, errors.New("get error type")
	}
	var iter = 0
	var ret []float64
	for {
		i := DecimalFloat(ByteToFloat(v[0+iter:2+iter]) * math.Pow(0.1, ByteToFloat(v[2+iter:4+iter])))
		ret = append(ret, i)
		if iter+4 == count {
			break
		} else {
			iter += 4
		}
	}
	return ret, nil
}

func DecimalFloat(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func TwoByteToFloatX1000(v []byte) (float64, error) {
	count := len(v)
	if count > 2 || count == 0 {
		return 0, errors.New("get error type")
	}
	var ret float64
	ret = ByteToFloat(v) / 1000
	return ret, nil
}

func T() {

}
