package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var port string = ":8181"

func main() {

	http.HandleFunc("/", mainPage)
	http.HandleFunc("/users", users)

	fmt.Println("Start server")

	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("End")
}

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	//user := User{FirstName: "Vasia", LastName: "Drek"}
	//js, err := json.Marshal(user)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//}
	w.Write([]byte(r.URL.Path))
	//fmt.Println(r.URL.Path)
}

func users(w http.ResponseWriter, r *http.Request) {
	userSlice := []User{User{"One", "Two"}, User{"Three", "Four"}}
	js, err := json.Marshal(userSlice)
	if err != nil {
		fmt.Println("Error:", err)
	}
	w.Write(js)
}
