package sensor

import (
	"net"
	"strings"
	"sync"
	"time"
)

type DeviceSession struct {
	readChan  chan []byte
	writeChan chan []byte
	stopChan  chan bool
	conn      net.Conn
	sync.Mutex
}

type SeaTurtle interface {
	GetDeviceData(d []byte)
}

var SessionsCollection sync.Map

/**
 * reg device to map
 */
func RegDeviceSession(conn net.Conn) DeviceSession {
	var s DeviceSession
	s.readChan = make(chan []byte)
	s.writeChan = make(chan []byte)
	s.stopChan = make(chan bool)
	s.conn = conn
	addr := strings.Split(conn.RemoteAddr().String(), ":")[0]
	SessionsCollection.Store(addr, s)
	return s
}

/**
 * return the device session which search for
 * if a non-existent, return {} and false
 */
func GetDeviceSession(addr string) (DeviceSession, bool) {
	if v, ok := SessionsCollection.Load(addr); ok {
		return v.(DeviceSession), true
	} else {
		return DeviceSession{}, false
	}
}

/**
 * send bytes to device
 * this Device must reg info last send message
 * the callback func will return data when the device send
 */
func (ds *DeviceSession) SendWord(data []byte, callback func(meta interface{}, data []byte)) {
	ds.Lock()
	ds.writeChan <- data
	for {
		select {
		case readData := <-ds.readChan:
			callback(8, readData)
			ds.Unlock()
			return
		}
	}
}

/**
 * read
 */
func (ds *DeviceSession) readConn() {
	for {
		data := make([]byte, 20, 20)
		if _, err := ds.conn.Read(data); err != nil {
			break
		}
		ds.readChan <- data
	}
	ds.stopChan <- true
}

/**
 * write
 */
func (ds *DeviceSession) writeConn() {
	for {
		data := <-ds.writeChan
		if _, err := ds.conn.Write(data); err != nil {
			break
		}
	}
	ds.stopChan <- true
}

// heart beating
func (ds *DeviceSession) HeartBeating(timeout int) {
	select {
	case _ = <-ds.readChan:
		print(ds.conn.RemoteAddr().String(), "keeping now")
		_ = ds.conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		break
	}
}
