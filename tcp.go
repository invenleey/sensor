package sensor

import (
	"fmt"
	"net"
	"time"
)

const (
	Network = "tcp"
	Address = ":6564"
)

// sensor status code testing
func testStatus()  {
	var i = 0
	for i = 1; i > 0; i++ {
		GetLocalDevicesInstance()
		v, _ := GetLocalSensor("7eb220dd-6127-58c7-8663-bf2f55371b78")
		fmt.Println(v.Status)
		time.Sleep(time.Second * 5)
	}
}

func RunDeviceTCP() {
	InitInfoMK()
	// go testStatus()
	listener, err := net.Listen(Network, Address)
	if err != nil {
		fmt.Println("[错误]", err)
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[错误]", err)
			return
		}
		go HandleProcessor(conn)
	}
}
