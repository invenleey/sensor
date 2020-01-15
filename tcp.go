package sensor

import (
	"fmt"
	"net"
	"time"
)

const (
	NETWORK = "tcp"
	ADDRESS = ":6564"
)

// sensor status code testing
func testStatus() {
	var i = 0
	for i = 1; i > 0; i++ {
		GetLocalDevicesInstance()
		v, _ := GetLocalSensor("7eb220dd-6127-58c7-8663-bf2f55371b78")
		fmt.Println(v.Status)
		time.Sleep(time.Second * 5)
	}
}

var listener net.Listener

/**
 * 关闭TCP
 */
func StopDeviceTCP() {
	listener.Close()
}

/**
 * 重启设备TCP
 */
func RestartDeviceTCP() {
	listener.Close()
	go RunDeviceTCP()

}

/**
 * 重启System
 */
func RestartTCPSystem() {
	listener.Close()
	s := GetDeviceSessions()
	s.Range(func(key, value interface{}) bool {
		h := value.(DeviceSession)
		h.stopChan <- true
		fmt.Println("[INFO] 已移除" + key.(string))
		return true
	})

	fmt.Println("[INFO] 重启TCP")
	go RunDeviceTCP()
}

func RunDeviceTCP() {
	// go testStatus()
	listener, _ = net.Listen(NETWORK, ADDRESS)

	// defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[FAIL] " + "退出TCP")
			return
		}
		go HandleProcessor(conn)
	}
}
