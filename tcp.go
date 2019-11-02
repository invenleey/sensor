package sensor

import (
	"fmt"
	"net"
)

func RunDeviceTCP() {
	listener, err := net.Listen("tcp", ":1080")
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("err =", err)
			return
		}
		go HandleProcessor(conn)
	}
}
