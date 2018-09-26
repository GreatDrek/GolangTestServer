package main

import (
	"fmt"
	"net"
)

func main() {
	buf := make([]byte, 1024)
	Conn, _ := net.DialUDP("udp", nil, &net.UDPAddr{IP: []byte{127, 0, 0, 1}, Port: 10001, Zone: ""})
	defer Conn.Close()
	go func() {
		for {
			n, addr, _ := Conn.ReadFromUDP(buf)
			fmt.Println(string(buf[0:n]), addr)
		}
	}()
	var input string
	for {
		fmt.Scanln(&input)
		if input == "close" {
			break
		}
		Conn.Write([]byte(input))
	}
}
