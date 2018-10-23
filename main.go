package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

var (
	addr *string
)

func main() {

	log.Println("Start server")

	port := os.Getenv("PORT")

	if port == "" {
		//log.Fatal("$PORT must be set")
		port = "5000"
	}

	addr = flag.String("addr", ":"+port, "http service address")

	flag.Parse()

	hub := newHub()

	go hub.run()

	//
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/bd", func(w http.ResponseWriter, r *http.Request) {
		creatDB()
	})
	//

	http.HandleFunc("/wss", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

//connStr := "user=postgres password=37352410 dbname=postgres sslmode=disable"
//db, err := sql.Open("postgres", connStr)
