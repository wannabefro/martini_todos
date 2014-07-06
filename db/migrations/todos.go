package main

import (
	_ "github.com/lib/pq"
	"database/sql"
	"github.com/coopernurse/gorp"
	"log"
)

func main() {
	dbmap := initDb()
	defer dbmap.Db.Close()
}

type Todo struct {
	Id					int64
	Created			int64
	Title				string
	Description	string
}

func initDb() *gorp.DbMap {
	db, err := sql.Open("postgres", "user=admin dbname=martinitodos sslmode=disable")
	checkErr(err, "sql.Open failed")

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	table := dbmap.AddTableWithName(Todo{}, "todos").SetKeys(true, "Id")
	table.ColMap("Title").SetNotNull(true)
	table.ColMap("Description").SetNotNull(true)

	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
