package sensor

import (
	"fmt"
	"testing"
)

func TestCheckCRC(t *testing.T) {
	dat := []byte{0x01, 0x02, 0x03, 0x04}
	m := []byte{0x2B, 0xA1}
	if ValidateCRC(dat, m) {
		fmt.Println("match")
	} else {
		fmt.Println("no match")
	}
}
