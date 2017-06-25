package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var (
	DBCon    *sql.DB
	connInfo = "postgres://teachme@127.0.0.1:5432/teachme?sslmode=disable"
)

func Initialize() (err error) {
	if DBCon != nil {
		err1 := DBCon.Close()
		if err1 != nil {
			return err1
		}
	}

	var err2 error
	DBCon, err2 = sql.Open("postgres", connInfo)
	if err2 != nil {
		return err2
	}
	return nil
}
