package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func Connect() *sqlx.DB {

	user := "ryouomoi-checker"
	password := "password"
	host := "localhost"
	port := "3306"
	database_name := "ryouomoi-checker-db"

	dbconf := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + database_name + "?charset=utf8mb4"
	db, err := sqlx.Open("mysql", dbconf)
	if err != nil {
		fmt.Println(err.Error())
	}
	return db
}
