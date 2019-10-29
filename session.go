package sensor

import "net"

type DeviceSession struct {
	readChan  chan []byte
	writeChan chan []byte
	stopChan  chan bool
	conn net.Conn
	callback func()
}

var SessionCollection map[string]DeviceSession

func SessionInit() {
	SessionCollection = make(map[string]DeviceSession)
}

func RegDevice(conn net.Conn) DeviceSession {
	var s DeviceSession
	s.readChan = make(chan []byte)
	s.writeChan = make(chan []byte)
	s.stopChan = make(chan bool)
	SessionCollection[conn.RemoteAddr().String()] = s
	return s
}
