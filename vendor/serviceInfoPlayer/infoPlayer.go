package serviceInfoPlayer

import (
	"database/sql"
	"log"
)

const NameTableInfoPlayer string = "infoplayers"

type InfoPlayer struct {
	Gold int `json:"gold"`
}

func (iP *InfoPlayer) LoadInfo(id int, db *sql.DB) error {
	row := db.QueryRow("SELECT id, gold FROM "+NameTableInfoPlayer+" WHERE id = $1", id)
	err := row.Scan(&id, &iP.Gold)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			_, err = db.Exec("INSERT INTO "+NameTableInfoPlayer+" (id, gold) VALUES ($1, $2)", id, 0)
			if err != nil {
				return err
			}
			iP.Gold = 0
			log.Println("Add info")
		} else {
			return err
		}
	} else {
		log.Println("Load info")
	}

	return err
}

func (iP *InfoPlayer) SaveInfo(id int, db *sql.DB) error {
	_, err := db.Exec("UPDATE "+NameTableInfoPlayer+" SET gold = $1 WHERE id = $2", iP.Gold, id)
	if err != nil {
		return err
	}
	return nil
}
