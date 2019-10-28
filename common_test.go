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