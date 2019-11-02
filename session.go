package sensor

import (
	"fmt"
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

type Ds interface {
	KillDevice()
	SendWord(data []byte, callback func(meta interface{}, data []byte))

	OpenReadTimeout()
	StopReadTimeout()

	ReadConn()
	WriteConn()
	HeartBeating(timeout int)
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
 * if a non-existent key, return {} and false
 */
func GetDeviceSession(addr string) (DeviceSession, bool) {
	if v, ok := SessionsCollection.Load(addr); ok {
		return v.(DeviceSession), true
	} else {
		return DeviceSession{}, false
	}
}

/**
 * delete key-value on sync hashMap
 */
func (ds *DeviceSession) KillDevice() {
	SessionsCollection.Delete(strings.Split(ds.conn.RemoteAddr().String(), ":")[0])
	fmt.Println(SessionsCollection)
}

/**
 * send bytes to device
 * this Device must reg info last send message
 * the callback func will return data when the device send
 */
func (ds *DeviceSession) SendWord(data []byte, callback func(meta interface{}, data []byte)) {
	ds.Lock()
	ds.writeChan <- data
	ds.OpenReadTimeout()
	for {
		select {
		case readData := <-ds.readChan:
			callback(8, readData)
			ds.StopReadTimeout()
			ds.Unlock()
			return
		}
	}
}

/**
 * read timeout open
 */
func (ds *DeviceSession) OpenReadTimeout() {
	if err := ds.conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		// error happened
		ds.stopChan <- true
	}
}

/**
 * read timeout stop
 */
func (ds *DeviceSession) StopReadTimeout() {
	// it is absoluteZeroYear = -292277022399
	if err := ds.conn.SetReadDeadline(time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)); err != nil {
		// error happened
		ds.stopChan <- true
	}
}

/**
 * read
 */
func (ds *DeviceSession) ReadConn() {
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
func (ds *DeviceSession) WriteConn() {
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
