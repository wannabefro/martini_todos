package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	_ "github.com/lib/pq"
	"github.com/coopernurse/gorp"
	"database/sql"
	"log"
	"./db/models"
)

func main() {
	db, err := sql.Open("postgres", "user=admin dbname=martinitodos sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	
	m := martini.Classic()
	m.Map(db)
	m.Use(martini.Static("assets"))

	m.Use(render.Renderer(render.Options {
		Layout: "layout",
	}))

	m.Get("/", func(r render.Render) {
		r.HTML(200, "welcome", map[string]interface{}{"greeting": "everybody"})
	})

	m.Get("/todos", func(r render.Render) {
		dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
		// defer dbmap.Db.Close()
		var todos []models.Todo
		_, err = dbmap.Select(&todos, "select * from todos order by created")
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(todos)
		r.HTML(200, "todos/index", todos)
	})

	m.Get("/todos/:id", func(params martini.Params, r render.Render) {
		dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
		// defer dbmap.Db.Close()
		var todo models.Todo
		err = dbmap.SelectOne(&todo, "select * from todos where id = :id",
		map[string]interface{} {"id": params["id"]})
		if err != nil {
			log.Println(err)
		}
		r.HTML(200, "todos/show", todo)
	})

	m.Run()
}
