package main

import (
	"database/sql"
	"flag"

	//"time"

	//"html/template"
	"log"
	"net/http"
	"os"

	"encoding/json"
	"io/ioutil"

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

	http.HandleFunc("/test/caroosel", CarooselPost)

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

func CarooselPost(w http.ResponseWriter, r *http.Request) {
	//log.Println(r.URL)
	if r.URL.Path != "/test/caroosel" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//http.ServeFile(w, r, "home.html")
	//log.Println(r.FormValue("id"))

	files, err := ioutil.ReadDir("./public/assets/img/photo_room/caroosel")
	if err != nil {
		log.Fatal(err)
	}

	imgCaroosel := make([]string, 0)

	for _, file := range files {
		if !file.IsDir() {
			imgCaroosel = append(imgCaroosel, file.Name())
		}
	}

	data, err := json.Marshal(imgCaroosel)
	if err != nil {
		log.Println(err)
	}

	w.Write(data)
}
