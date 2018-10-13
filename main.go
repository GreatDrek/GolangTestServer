package main

import (
	//"myLib"
	"connect"
	"log"

	//"github.com/gorilla/websocket"

	//"time"

	"fmt"
	"net/http"
)

func main() {

	fmt.Println("Start server")

	port := "5000"

	go httpServer(port)

	fmt.Scanln()
}

func httpServer(port string) {
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/status", statusServer)

	http.HandleFunc("/test", test)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.Write([]byte("error"))
	} else {
		if r.URL.Path != "/favicon.ico" {
			b := []byte(r.URL.Path)
			connect.AddUser(string(b[1:]))
			w.Write(b[1:])
		}
	}
}

func statusServer(w http.ResponseWriter, r *http.Request) {
	s := ""
	for key, value := range connect.ReturnAllConnect() {
		s += key + ": " + value + "\n"
	}
	w.Write([]byte(s))
}

func test(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://127.0.0.1:5000/status")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	for true {
		bs := make([]byte, 1014)
		n, err := resp.Body.Read(bs)

		fmt.Println(string(bs[:n]))

		if n == 0 || err != nil {
			break
		}
	}
}
