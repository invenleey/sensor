package sensor

import "net"

type DeviceSession struct {
	readChan  chan []byte
	writeChan chan []byte
	stopChan  chan bool
	conn net.Conn
}

var SessionCollection map[string]DeviceSession

/**
 * sessions init
 */
func SessionsInit() {
	SessionCollection = make(map[string]DeviceSession)
}

/**
 * reg device to map
 */
func RegDevice(conn net.Conn) DeviceSession {
	var s DeviceSession
	s.readChan = make(chan []byte)
	s.writeChan = make(chan []byte)
	s.stopChan = make(chan bool)
	SessionCollection[conn.RemoteAddr().String()] = s
	return s
}

/**
 * send bytes to device
 * this Device must reg info last send message
 */
func SendBytes(addr string, b []byte) error {
	SessionCollection[addr].writeChan <- b
	return nil
}

func Send(addr string, b []byte, callback func()) error {
	SessionCollection[addr].writeChan <- b
	print("send inner")
	return nil
}