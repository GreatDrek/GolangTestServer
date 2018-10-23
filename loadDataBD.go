package main

import (
	//"crypto/sha256"

	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

type bdInfo struct {
	id    string
	email string
	key   []byte
	salt  []byte
}

func checkUser(logData logginDataClient) (*bdInfo, error) {
	var err error
	var infoClient *bdInfo

	//connStr := "user=postgres password=37352410 dbname=postgres sslmode=disable"
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	//db, err = sql.Open("postgres", connStr)
	if err != nil {
		//log.Fatalf("Error opening database: %q", err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, email, key, salt FROM idusers WHERE email = $1", logData.Email)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var id string
	var email string
	var key []byte
	var salt []byte

	for rows.Next() {
		if err = rows.Scan(&id, &email, &key, &salt); err != nil {
			return nil, err
		} else {
			infoClient = &bdInfo{id: id, email: email, key: key, salt: salt}
		}
	}

	return infoClient, err
}

func addUser(infoClient *bdInfo) error {
	var err error
	//connStr := "user=postgres password=37352410 dbname=postgres sslmode=disable"
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	//db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	log.Println(infoClient)
	//
	// IDENTITY чекнуть!!!
	//
	if _, err := db.Exec("INSERT INTO idusers VALUES ($1, $2, $3, $4)", "0", infoClient.email, infoClient.key, infoClient.salt); err != nil {
		log.Println("test", err)
		return err
	}
	return err
}

func updateUser(infoClient *bdInfo) error {
	var err error
	//connStr := "user=postgres password=37352410 dbname=postgres sslmode=disable"
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	//db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	log.Println(infoClient)
	//
	// IDENTITY чекнуть!!!
	//
	//if _, err := db.Exec("INSERT INTO idusers VALUES ($1, $2, $3, $4)", "0", infoClient.email, infoClient.key, infoClient.salt); err != nil {
	//	log.Println("test", err)
	//	return err
	//}

	_, err = db.Exec("UPDATE idusers SET key = $1, salt = $2 WHERE email = $3", infoClient.key, infoClient.salt, infoClient.email)
	if err != nil {
		return err
	}

	return err
}

func creatDB() {
	//connStr := "user=postgres password=37352410 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	//db, err := sql.Open("postgres", connStr)
	if err != nil {
		return
	}
	defer db.Close()

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS idusers (id text, email text, key bytea, salt bytea)"); err != nil {
		return
	}
}
