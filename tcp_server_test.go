package sensor

import "testing"

func TestOpenKit(t *testing.T) {
	go RunDeviceTCP()
	OpenKit()
}