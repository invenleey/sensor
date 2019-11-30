package sensor

import (
	"errors"
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
	interfaceDevice
}

type interfaceDevice interface {
	KillDevice()
	SendWord(data []byte, callback func(dm DeviceMeta, data []byte) (ReadResult, error)) (ReadResult, error)

	OpenReadTimeout()
	StopReadTimeout()

	ReadConn()
	WriteConn()
	HeartBeating(timeout int)

	GetResultInstance(meta DeviceMeta) (ReadResult, error)
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
 * show all device node ip in this -> lan
 * @return a string array include all node server ip
 */
func ShowNodeIPs() []string {
	var ret []string
	SessionsCollection.Range(func(key, value interface{}) bool {
		ret = append(ret, key.(string))
		return true
	})
	return ret
}

/**
 * delete key-value on sync hashMap
 */
func (ds *DeviceSession) KillDevice() {
	SessionsCollection.Delete(strings.Split(ds.conn.RemoteAddr().String(), ":")[0])
}

type DeviceMeta struct {
	Addr     byte
	FuncCode byte
}

//type interfaceSplit interface {
//	ByteSplit(src []byte) (DeviceMeta, []byte, error)
//	GetData() []byte
//}

/**
 * send bytes to device
 * this Device must reg info last send message
 * the callback func will return data when the device send
 */
func (ds *DeviceSession) SendWord(data []byte, callback func(dm DeviceMeta, data []byte) (ReadResult, error)) (ReadResult, error) {
	ds.Lock()
	ds.writeChan <- data
	ds.OpenReadTimeout()
	for {
		select {
		case readData := <-ds.readChan:
			ds.StopReadTimeout()
			dm, md, err := SplitMeasure(readData)
			var rs ReadResult
			if err != nil {
				fmt.Println("error data")
			} else {
				rs, err = callback(dm, md)
			}
			ds.Unlock()
			return rs, err
		}
	}
}

/**
 * read timeout open
 */
func (ds *DeviceSession) OpenReadTimeout() {
	if err := ds.conn.SetReadDeadline(time.Now().Add(15 * time.Second)); err != nil {
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

		data := make([]byte, 20)
		// var data []byte
		n, err := ds.conn.Read(data)
		if err != nil {
			break
		}
		ds.readChan <- data[:n]
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

func (ds *DeviceSession) MeasureRequest(rData []byte, itemsName []string) (ReadResult, error) {
	p, err := ds.SendWord(rData, func(meta DeviceMeta, data []byte) (ReadResult, error) {
		p, err := ds.GetResultInstance(meta)
		if err != nil {
			return ReadResult{}, errors.New("meta build error")
		}
		if err = p.DecodeStandardFourByte2Float(data, itemsName); err != nil {
			return ReadResult{}, errors.New("decode build error")
		} else {
			return p, nil
		}
	})
	if err != nil {
		return ReadResult{}, err
	} else {
		return p, nil
	}

	//go b.SendWord([]byte{0x06, 0x03, 0x10, 0x06, 0x00, 0x01, 0x61, 0x7C}, func(meta DeviceMeta, data []byte) {
	//	p, err := b.GetReadResultInstance(meta)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	if err = p.DecodeSlope(data, "斜率校准值"); err != nil {
	//		fmt.Println(err)
	//	} else {
	//		fmt.Println(p)
	//	}
	//})

	//go b.SendWord([]byte{0x06, 0x06, 0x20, 0x02, 0x00, 0x01, 0xE3, 0xBD}, func(meta DeviceMeta, data []byte) {
	//	p, err := b.GetResultInstance(meta)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	if err = p.DecodeOrder(data); err != nil {
	//		fmt.Println(err)
	//	} else {
	//		fmt.Println(p)
	//	}
	//})
}
