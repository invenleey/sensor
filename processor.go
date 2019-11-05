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
	go b.SendWord([]byte{0x06, 0x03, 0x00, 0x00, 0x00, 0x04, 0x45, 0xBE}, func(meta DeviceMeta, data []byte) {
		p := b.GetMeasureResultInstance()
		_ = p.DecodeMeasureByte(meta, data, 0x03, []string{"测量值", "温度"})
		fmt.Println(p)
	})

	go b.SendWord([]byte{0x01}, func(meta DeviceMeta, data []byte) {
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
