package main

import (
	"database/sql"
	"flag"

	//"time"

	//"html/template"
	"log"
	"net/http"
	"os"

	//"path"
	"serviceAutorization"
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

	hub := serviceConnection.NewHub()
	hub.Db = db

	go hub.Run()

	http.HandleFunc("/u", serveHome)

	fs1 := http.FileServer(http.Dir("./public/"))
	http.Handle("/test/", http.StripPrefix("/test/", fs1))

	http.HandleFunc("/wss", func(w http.ResponseWriter, r *http.Request) {
		serviceConnection.ServeWs(hub, w, r, &Client{})
	})

	http.HandleFunc("/db", func(w http.ResponseWriter, r *http.Request) {
		serviceAutorization.CreatDB(db)
	})

	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/u" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}
