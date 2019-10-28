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

func TestSplitBody(t *testing.T) {
	rs := []byte{0x01, 0x06, 0x20, 0x02, 0x00, 0x01, 0xe3, 0xbd}

}
