package sensor

import (
	"fmt"
	"net"
	"time"
)

const (
	NETWORK = "tcp"
	ADDRESS = ":6564"
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

/**
 * 对传感器的服务重启
 * 尽量避免重启解决问题, 需要重启说明程序不够完美
 */

func RunDeviceTCP() {
	InitInfoMK()
	go testStatus()
	listener, err := net.Listen(NETWORK, ADDRESS)
	if err != nil {
		fmt.Println("[FAIL]", err)
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[FAIL]", err)
			return
		}
		go HandleProcessor(conn)
	}
}
