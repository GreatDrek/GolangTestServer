package main

import (
	"fmt"
	"net"
	"time"
)

type myPackage struct {
	idServer        string
	idClient        string
	idMyPackage     int32
	idServerPackage int32
}

func main() {
	myPackage := myPackage{idServer: "fH8", idClient: "0000"}

	buf := make([]byte, 1024)
	Conn, _ := net.DialUDP("udp", nil, &net.UDPAddr{IP: []byte{127, 0, 0, 1}, Port: 10001, Zone: ""})
	defer Conn.Close()
	go func() {
		for {
			n, addr, _ := Conn.ReadFromUDP(buf)
			fmt.Println(string(buf[0:n]), addr)
		}
	}()
	go func() {
		for {
			time.Sleep(time.Millisecond * 35)
			Conn.Write([]byte(myPackage.idServer + myPackage.idClient))
		}
	}()
	var input string
	for {
		fmt.Scanln(&input)
		if input == "close" {
			break
		}
		//Conn.Write([]byte("fH8" + input))
	}
}
