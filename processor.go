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
	fmt.Println(conn.RemoteAddr(), "已经连接")
	defer conn.Close()
	// session
	b := RegDevice(conn)

	go readConn(conn, b.readChan, b.stopChan)
	go writeConn(conn, b.writeChan, b.stopChan)
	// go HeartBeating(conn, readChan, 20)
	for {
		select {
		case readStr := <-b.readChan:
			// callback
			getData(readStr)
			// b.writeChan <- readStr
		case stop := <-b.stopChan:
			// 弹出
			if stop {
				fmt.Println(conn.RemoteAddr(), "已经断开连接")
				break
			}
		}
	}
}

func getData(msg []byte) {
	fmt.Println("某一协程从服务器接收到: ", msg)
}

func readConn(conn net.Conn, readChan chan<- []byte, stopChan chan<- bool) {
	for {
		data := make([]byte, 20, 20)
		_, err := conn.Read(data)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println("Received:", data)
		readChan <- data
	}
	stopChan <- true
}
func writeConn(conn net.Conn, writeChan <-chan []byte, stopChan chan<- bool) {
	// test
	k := []byte{0x06, 0x03, 0x00, 0x00, 0x00, 0x04, 0x45, 0xBE}
	_, _ = conn.Write(k)
	for {
		strData := <-writeChan
		_, err := conn.Write(strData)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println("Send:", strData)
	}
	stopChan <- true
}

// 心跳检测
func HeartBeating(conn net.Conn, readerChannel chan []byte, timeout int) {
	select {
	case _ = <-readerChannel:
		print(conn.RemoteAddr().String(), "get message, keeping heartbeating...")
		_ = conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		break
	}
}
