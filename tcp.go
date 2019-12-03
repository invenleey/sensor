package sensor

import (
	"fmt"
	"net"
)

const (
	Network = "tcp"
	Address = ":6564"
)

func RunDeviceTCP() {
	InitInfoMK()
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
