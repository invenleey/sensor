package sensor

import (
	"bytes"
	"encoding/binary"
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
