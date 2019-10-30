package sensor

import (
	"fmt"
	"net"
	"time"
)





func HandleProcessor(conn net.Conn) {
	fmt.Println("[连接]", conn.RemoteAddr())
	defer conn.Close()
	// session
	b := RegDevice(conn)

	go readConn(conn, b.readChan, b.stopChan)
	go writeConn(conn, b.writeChan, b.stopChan)
	// go HeartBeating(conn, readChan, 20)

	// testing
	go SendWord(conn.RemoteAddr().String(), []byte{0x00}, func(meta interface{}, data []byte) {
		fmt.Println(data)
	})

	for {
		select {
		// abandon function
		//case readStr := <-b.readChan:
		//	getData(readStr)
		case stop := <-b.stopChan:
			// 弹出
			if stop {
				fmt.Println("[断开]", conn.RemoteAddr())
				break
			}
		}
	}
}

//func SendMeasure(addr string, data []byte, callback func(meta interface{}, data []byte)) {
//	SessionCollection[addr].writeChan <- data
//	for {
//		select {
//		case readData := <-SessionCollection[addr].readChan:
//			callback(9, readData)
//			return
//		}
//	}
//}

func SendWord(addr string, data []byte, callback func(meta interface{}, data []byte)) {
	SessionCollection[addr].writeChan <- data
	for {
		select {
		case readData := <-SessionCollection[addr].readChan:
			callback(8, readData)
			return
		}
	}
}

func readConn(conn net.Conn, readChan chan<- []byte, stopChan chan<- bool) {
	for {
		data := make([]byte, 20, 20)
		if _, err := conn.Read(data); err != nil {
			break
		}
		readChan <- data
	}
	stopChan <- true
}

func writeConn(conn net.Conn, writeChan <-chan []byte, stopChan chan<- bool) {
	for {
		data := <-writeChan
		if _, err := conn.Write(data); err != nil {
			break
		}
	}
	stopChan <- true
}

// 心跳检测
func HeartBeating(conn net.Conn, readerChannel chan []byte, timeout int) {
	select {
	case _ = <-readerChannel:
		print(conn.RemoteAddr().String(), "keeping now")
		_ = conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		break
	}
}
