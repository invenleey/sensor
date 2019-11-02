package sensor

import (
	"fmt"
	"net"
)

func HandleProcessor(conn net.Conn) {
	fmt.Println("[连接]", conn.RemoteAddr())
	defer conn.Close()
	// session
	b := RegDeviceSession(conn)
	go b.ReadConn()
	go b.WriteConn()
	// go b.HeartBeating(20)

	// testing
	go b.SendWord([]byte{0x00}, func(meta interface{}, data []byte) {
		fmt.Println(data)
	})

	go b.SendWord([]byte{0x01}, func(meta interface{}, data []byte) {
		fmt.Println(data)
	})

	for {
		select {
		case stop := <-b.stopChan:
			// pick out
			if stop {
				fmt.Println("[断开]", conn.RemoteAddr())
				b.KillDevice()
				break
			}
		}
	}
}
