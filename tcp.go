package sensor

import (
	"bufio"
	"fmt"
	"net"
)

func RunDeviceTCP() {
	listener, err := net.Listen("tcp", ":1080")
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("err =", err)
			return
		}
		// go HandleConn(conn)
		// go process(conn)
		go HandleD(conn)
	}
}

func process(conn net.Conn) {
	defer conn.Close()
	processor := &Processor{
		Conn: conn,
	}
	processor.process2()
}

func HandleConn(conn net.Conn) {
	defer conn.Close()
	addr := conn.RemoteAddr().String()
	fmt.Println(addr, "已连接到服务器")
	reader := bufio.NewReader(conn)
	k := []byte{0x06, 0x03, 0x00, 0x00, 0x00, 0x04, 0x45, 0xBE}
	_, _ = conn.Write(k)
	for {
		if l, err := reader.ReadByte(); err != nil {
			fmt.Println(addr, "已断开连接")
			return
		} else {
			fmt.Println(l)
		}
	}
}

