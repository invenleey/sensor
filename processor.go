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

// Device Info and the processor detail context
type Processor struct {
	// TCP Connection
	Conn net.Conn
	// The ModBus Address which device setup
	// Addr uint8
	// Operation Type
	Func uint8
	S    Sensor
}

var request []byte

func RequestBuild() {

}

func GetSensorAddr() {

}

func (p *Processor) process2() {
	defer p.Conn.Close()
	addr := p.Conn.RemoteAddr().String()
	fmt.Println(addr, "已连接到服务器")
	// reader := bufio.NewReader(p.Conn)
	k := []byte{0x06, 0x03, 0x00, 0x00, 0x00, 0x04, 0x45, 0xBE}
	_, _ = p.Conn.Write(k)
	var bk []byte
	for {
		if _, err := p.Conn.Read(bk); err != nil {
			fmt.Println(addr, "已断开连接")
			return
		} else {
			fmt.Println(bk)
		}
	}
}

func HandleD(conn net.Conn) {
	fmt.Println(conn.RemoteAddr(), "已经连接")
	defer conn.Close()
	readChan := make(chan []byte)
	writeChan := make(chan []byte)
	stopChan := make(chan bool)
	go readConn(conn, readChan, stopChan)
	go writeConn(conn, writeChan, stopChan)
	go HeartBeating(conn, readChan, 20)
	for {
		select {
		case readStr := <-readChan:
			getData(readStr)
			writeChan <- readStr
		case stop := <-stopChan:
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

// 心跳检测，根据GravelChannel判断Client是否在设定时间内发来信息
func HeartBeating(conn net.Conn, readerChannel chan []byte, timeout int) {
	select {
	case _ = <-readerChannel:
		print(conn.RemoteAddr().String(), "get message, keeping heartbeating...")
		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		break
	}
}
