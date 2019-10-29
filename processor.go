package sensor

import (
	"fmt"
	"net"
	"time"
)

// Function Code Type
// Read 0x03
// Write 0x06
var ReadFunc = []byte{0x03}
var WriteFunc = []byte{0x06}

// Register Address
var RRegMeasure = []byte{0x00, 0x00}
var WRegOxygen = []byte{0x10, 0x04}
var WRegZero = []byte{010, 0x00}
var WRegTilt = []byte{0x10, 0x04}

var RRegZero = []byte{010, 0x06}
var RRegTilt = []byte{0x10, 0x08}

var ARegAddr = []byte{0x20, 0x02}
var WRegFactory = []byte{0x20, 0x20}

var request []byte

func HandleProcessor(conn net.Conn) {
	fmt.Println("[连接]", conn.RemoteAddr())
	defer conn.Close()
	// session
	b := RegDevice(conn)

	go readConn(conn, b.readChan, b.stopChan)
	go writeConn(conn, b.writeChan, b.stopChan)
	// go HeartBeating(conn, readChan, 20)
	_ = SendBytes(conn.RemoteAddr().String(), []byte{0x01})
	for {
		select {
		case readStr := <-b.readChan:
			// callback
			getData(readStr)
			// b.writeChan <- readStr
		case stop := <-b.stopChan:
			// 弹出
			if stop {
				fmt.Println("[断开]", conn.RemoteAddr())
				break
			}
		}
	}
}

func getData(msg []byte) {
	fmt.Println("Got: ", msg)
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
