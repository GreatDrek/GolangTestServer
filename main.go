package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"serviceConnection"

	_ "github.com/lib/pq"
)

var (
	addr *string
	db   *sql.DB
)

func main() {
	//connStr := "user=postgres password=37352410 dbname=postgres sslmode=disable"
	//db, err = sql.Open("postgres", connStr)

	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("BD don't connect", err)
	}
	defer db.Close()

	log.Println("Start server")

	port := os.Getenv("PORT")

	if port == "" {
		//log.Fatal("$PORT must be set")
		port = "5000"
	}

	addr = flag.String("addr", ":"+port, "http service address")

	flag.Parse()

	//hub := newHub()
	hub := serviceConnection.NewHub()

	go hub.Run()

	http.HandleFunc("/", serveHome)

	http.HandleFunc("/wss", func(w http.ResponseWriter, r *http.Request) {
		serviceConnection.ServeWs(hub, w, r, &Client{})
	})

	err = http.ListenAndServe(*addr, nil)
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
