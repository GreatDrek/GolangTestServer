package main

import (
	"fmt"
	"net"
)

type conn struct {
	packgID         string //3 numbers fH8
	indentification string //4 numbers
	addr            *net.UDPAddr
}

func main() {
	fmt.Println("Start server")

	go Connect()

	fmt.Scanln()

	fmt.Println("Stop server")
}

func Connect() {
	ServerConn, _ := net.ListenUDP("udp", &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: 10001, Zone: ""})
	defer ServerConn.Close()
	buf := make([]byte, 1024)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)

		if err != nil {
			fmt.Println("Error 0:", err)
			continue
		}

		if n >= 7 {
			if string(buf[0:3]) == "fH8" {
				fmt.Println(addr, string(buf[0:n]))
			}
		}
	}
}
