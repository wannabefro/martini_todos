package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	_ "github.com/lib/pq"
	// "github.com/coopernurse/gorp"
	"database/sql"
	"log"
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

	m.Run()
}
