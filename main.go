package main

import (
	"fmt"
	//"io"
	"bufio"
	"net"
)

func main() {
	listener, _ := net.Listen("tcp", ":5000")
	//allConn := make([]net.Conn, 0, 5)

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Dont connect:", err)
			conn.Close()
			continue
		}

		//allConn = append(allConn, conn)

		fmt.Println("Connect")

		bufReader := bufio.NewReader(conn)
		fmt.Println("Start reading")

		go func(conn net.Conn) {
			defer conn.Close()
			for {
				rbyte, _, err := bufReader.ReadLine()

				if err != nil {
					fmt.Println("Buf error:", err)
					break
				}

				s := string(rbyte)

				fmt.Println(s)

				//for _, val := range allConn {
				//	val.Write([]byte(s))
				//}
			}
		}(conn)
	}
	//listener.Close()
	//fmt.Println("Stop server")
}
