package main

import (
	"fmt"
	"net"
)

type Client struct {
	index    byte
	addr     *net.UDPAddr
	nickName string
}

func main() {
	ServerConn, _ := net.ListenUDP("udp", &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: 10001, Zone: ""})
	defer ServerConn.Close()
	buf := make([]byte, 1024)

	allClient := make(map[byte]Client)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)

		if err != nil {
			fmt.Println("Error 0:", err)
		}

		allClient[buf[0]] = Client{index: buf[0], addr: addr, nickName: string(buf[1:n])}

		//s := string(buf[1:n])
	}

	//	allAddr := make([]*net.UDPAddr, 0, 10)
	//	for {
	//		n, addr, _ := ServerConn.ReadFromUDP(buf)

	//		if len(allAddr) == 0 {
	//			allAddr = append(allAddr, addr)
	//		} else {
	//			for idx, val0 := range allAddr {
	//				if idx == len(allAddr)-1 {
	//					if val0.Port != addr.Port {
	//						allAddr = append(allAddr, addr)
	//					} else {
	//						break
	//					}
	//				} else {
	//					if val0.Port == addr.Port {
	//						break
	//					}
	//				}
	//			}
	//		}

	//		for _, val := range allAddr {
	//			//if val, ok := allAddr[key]; ok{
	//			//fmt.Println(val)
	//			ServerConn.WriteToUDP([]byte(string(buf[0:n])), val)
	//			//}
	//			//ServerConn.WriteToUDP([]byte("TEST"), val)
	//			fmt.Println(allAddr)
	//		}
	//		fmt.Println("Received ", string(buf[0:n]), " from ", addr)
	//		fmt.Println(len(allAddr))
	//}
}
