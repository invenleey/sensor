package sensor

import (
	"fmt"
	"net"
	"strings"
)

func HandleProcessor(conn net.Conn) {
	defer conn.Close()
	b := RegDeviceSession(conn)

	go b.ReadConn()
	go b.WriteConn()
	// go b.HeartBeating(20)

	dtuIpv4 := strings.Split(conn.RemoteAddr().String(), ":")[0]

	fmt.Println("[CONN]", conn.RemoteAddr())

	// setup time wheel
	ch := TaskSetup(dtuIpv4)

	// fmt.Println("already connected:", ShowNodeIPs())

	// testing
	//
	//go b.SendWord([]byte{0x01}, func(meta DeviceMeta, data []byte) {
	//	fmt.Println(data)
	//})

	//for {
		select {
		case stop := <-b.stopChan:
			// pick out
			if stop {
				fmt.Println("[DISC]", conn.RemoteAddr())
				// 资源释放过程

				// 移除session
				b.ReleaseDevice()

				// 移除task
				b.ReleaseTask()

				// 关闭任务通道
				close(ch)

				break
			}
		}
	//}
}
