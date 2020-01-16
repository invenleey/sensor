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
	// dtu设备释放
	ReleaseDevice()
	// 定时任务释放
	ReleaseTask()

	// 写1
	SendWord(data []byte, callback func(dm DeviceMeta, data []byte) (ReadResult, error)) (ReadResult, error)
	// 写2
	SendToSensor(requestData []byte) ([]byte, error)

	// TCP超时
	OpenReadTimeout()
	StopReadTimeout()

	// 读写
	ReadConn()
	WriteConn()

	// 心跳
	HeartBeating(timeout int)

	// 内容, 不应该在这
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
 * if a non-existent key, return err
 */
func GetDeviceSession(addr string) (DeviceSession, error) {
	if v, ok := SessionsCollection.Load(addr); ok {
		return v.(DeviceSession), nil
	} else {
		return DeviceSession{}, errors.New("not found session")
	}
}

func GetDeviceSessions() sync.Map {
	return SessionsCollection
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
 * 释放Map中session
 */
func (ds *DeviceSession) ReleaseDevice() {
	SessionsCollection.Delete(strings.Split(ds.conn.RemoteAddr().String(), ":")[0])
}

/**
 * 释放Map中task
 */
func (ds *DeviceSession) ReleaseTask() {
	// 移除任务
	for _, v := range GetLocalDevicesInstance().GetLocalSensorList(strings.Split(ds.conn.RemoteAddr().String(), ":")[0]) {
		if err := v.RemoveTask(); err != nil {
			fmt.Println("[WARN] 释放任务过程出现错误 ", err)
		} else {
			fmt.Println("[INFO] 释放传感器任务资源 " + v.SensorID)
		}
	}
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
 * 简单发送
 * 超时返回
 *
 */
func (ds *DeviceSession) SendToSensor(requestData []byte) ([]byte, error) {
	ds.writeChan <- requestData
	for {
		select {
		case readData := <- ds.readChan:
			return readData, nil
		case <-time.After(10 * time.Second):
			return nil, errors.New("connected timeout")
		}
	}
}

/**
 * 向ds发送特定的指令(需要包含crc纠错), 等待回应
 * @param data 请求体
 * @callback 自定义的回调处理
 * @param timeout 超时channel处理
 */
func (ds *DeviceSession) SendWord(data []byte, callback func(dm DeviceMeta, data []byte) (ReadResult, error)) (ReadResult, error) {
	ds.writeChan <- data
	for {
		select {
		case readData := <-ds.readChan:
			// 检测数据
			dm, md, err := SplitAndValidate(readData)
			var rs ReadResult
			if err != nil {
				fmt.Println("error data")
			} else {
				// 回调的自定义处理
				rs, err = callback(dm, md)
			}
			return rs, err
		case <-time.After(10 * time.Second):
			// 超时处理, 识别为不存在的传感器, 即失去物理连接的
			// 是否考虑多次才出现
			fmt.Println("[WARN] 传感器连接超时")
			return ReadResult{}, errors.New("sensor timeout")
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
