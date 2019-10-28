package sensor

import (
	"fmt"
	"testing"
)

func TestToLittleEndian(t *testing.T) {
	dat := []byte{0x01, 0x02, 0x03, 0x04}
	checksum := CheckSum(dat)
	fmt.Printf("check sum:%X \n", ToLittleEndian(checksum))
}

func TestToBigEndian(t *testing.T) {
	dat := []byte{0x01, 0x02, 0x03, 0x04}
	checksum := CheckSum(dat)
	fmt.Printf("check sum:%X \n", ToBigEndian(checksum))
}

func TestComposeBody(t *testing.T) {
	addr := []byte{0x06, 0x03}
	funcdata := []byte{0x00, 0x00}
	data := []byte{0x00, 0x04}
	rs := ComposeBody(addr, funcdata, data)
	fmt.Println(rs)
}

func TestSplitConfig(t *testing.T) {
	rs := []byte{0x01, 0x06, 0x20, 0x02, 0x00, 0x01, 0xE2, 0x0A}
	a, b, c, e := SplitConfig(rs)
	println(a, b, c, e)
	if e != nil {
		t.Fail()
	}
}

func TestSplitMeasure(t *testing.T) {
	rs := []byte{0x06, 0x03, 0x08, 0x04, 0x7F, 0x00, 0x02, 0x01, 0x1E, 0x00, 0x01, 0xD9, 0x6D}
	a, b, e := SplitMeasure(rs)
	println(a, b, e)
	if e != nil {
		t.Fail()
	}
}
