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
	//go b.SendWord([]byte{0x06, 0x03, 0x00, 0x00, 0x00, 0x04, 0x45, 0xBE}, func(meta DeviceMeta, data []byte) {
	//	p, err := b.GetReadResultInstance(meta)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	if err = p.DecodeMeasureByte(data, []string{"测量值", "温度"}); err != nil {
	//		fmt.Println(err)
	//	} else {
	//		fmt.Println(p)
	//	}
	//})

	go b.SendWord([]byte{0x06, 0x03, 0x10, 0x06, 0x00, 0x01, 0x61, 0x7C}, func(meta DeviceMeta, data []byte) {
		p, err := b.GetReadResultInstance(meta)
		if err != nil {
			fmt.Println(err)
		}
		if err = p.DecodeSlope(data, "零点校准值"); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(p)
		}
	})
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
