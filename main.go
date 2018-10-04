package main

import (
	"fmt"
	"net"
	//"time"
)

var s string

func main() {
	listener, err := net.Listen("tcp", ":4545")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}
		go handleConnection(conn) // запускаем горутину для обработки запроса
	}
}

// обработка подключения
func handleConnection(conn net.Conn) {
	defer conn.Close()
	//	go func() {
	//		for {
	//			_, err := conn.Write([]byte(s))
	//			if err != nil {
	//				fmt.Println("Close connect")
	//				break
	//			}
	//			//conn.Write([]byte(s))
	//			time.Sleep(time.Second * 10)
	//			fmt.Println(conn)
	//		}
	//	}()
	for {
		// считываем полученные в запросе данные
		input := make([]byte, (1024 * 4))
		n, err := conn.Read(input)
		if n == 0 || err != nil {
			fmt.Println("Read error:", err)
			break
		}
		source := string(input[0:n])

		s = source

		// выводим на консоль сервера диагностическую информацию
		fmt.Println(source)
		// отправляем данные клиенту
		//conn.Write([]byte(s))
	}

	fmt.Println("Exit")
}
