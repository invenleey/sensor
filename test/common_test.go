package test

import (
	"fmt"
	"sensor"
	"testing"
)

func TestToLittleEndian(t *testing.T) {
	dat := []byte{0x01, 0x02, 0x03, 0x04}
	checksum := sensor.CheckSum(dat)
	fmt.Printf("check sum:%X \n", sensor.ToLittleEndian(checksum))
}

func TestToBigEndian(t *testing.T) {
	dat := []byte{0x01, 0x02, 0x03, 0x04}
	checksum := sensor.CheckSum(dat)
	fmt.Printf("check sum:%X \n", sensor.ToBigEndian(checksum))
}

func TestComposeBody(t *testing.T) {
	addr := []byte{0x06, 0x03}
	funcdata := []byte{0x00, 0x00}
	data := []byte{0x00, 0x04}
	rs := sensor.ComposeBody(addr, funcdata, data)
	fmt.Println(rs)
}

func TestSplitConfig(t *testing.T) {
	rs := []byte{0x01, 0x06, 0x20, 0x02, 0x00, 0x01, 0xE2, 0x0A}
	a, b, c, e := sensor.SplitConfig(rs)
	println(a, b, c, e)
	if e != nil {
		t.Fail()
	}
}

func TestSplitMeasure(t *testing.T) {
	rs := []byte{0x06, 0x03, 0x08, 0x04, 0x7F, 0x00, 0x02, 0x01, 0x1E, 0x00, 0x01, 0xD9, 0x6D}
	a, b, e := sensor.SplitAndValidate(rs)
	println(a.FuncCode, b, e)
	if e != nil {
		t.Fail()
	}
}

func TestFourByteToFloat(t *testing.T) {
	// 4 bytes output
	v := []byte{0x01, 0x02, 0x00, 0x02}
	fmt.Println(sensor.FourByteToFloat(v))
	v = []byte{0x00, 0xB0, 0x00, 0x01}
	fmt.Println(sensor.FourByteToFloat(v))

	// 8 bytes output
	v = []byte{0x01, 0x02, 0x00, 0x02, 0x00, 0xB0, 0x00, 0x01}
	fmt.Println(sensor.FourByteToFloat(v))

	//16 byte output
	v = []byte{0x01, 0x02, 0x00, 0x02, 0x00, 0xB0, 0x00, 0x01, 0x01, 0x02, 0x00, 0x02, 0x00, 0xB0, 0x00, 0x01}
	fmt.Println(sensor.FourByteToFloat(v))

	// error output type(error byte count)
	v = []byte{0x01, 0x02, 0x00}
	_, err := sensor.FourByteToFloat(v)
	if err == nil {
		t.Fail()
	}
	// error output type(nil byte input)
	v = []byte{}
	_, err = sensor.FourByteToFloat(v)
	if err == nil {
		t.Fail()
	}
}

func TestByteToFloat(t *testing.T) {
	v := []byte{0x1, 0x2}
	fmt.Println(sensor.ByteToFloat(v))
}

//func TestReadDeviceList(t *testing.T) {
//	entity := LoadConfig("conf.json")
//	fmt.Println(entity)
//}

func TestDumpConfig(t *testing.T) {
	_ = sensor.GetLocalDevicesInstance().DumpConfig()
}