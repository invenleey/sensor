package sensor

import (
	"fmt"
	"net"
	"strings"
)

func HandleProcessor(conn net.Conn) {
	fmt.Println("[连接]", conn.RemoteAddr())
	defer conn.Close()
	// session
	b := RegDeviceSession(conn)
	// setup time wheel
	TaskSetup(strings.Split(conn.RemoteAddr().String(), ":")[0])

	go b.ReadConn()
	go b.WriteConn()
	// go b.HeartBeating(20)

	// fmt.Println("already connected:", ShowNodeIPs())

	// testing
	//
	//go b.SendWord([]byte{0x01}, func(meta DeviceMeta, data []byte) {
	//	fmt.Println(data)
	//})

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
