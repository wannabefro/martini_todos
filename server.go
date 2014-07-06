package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/binding"
	_ "github.com/lib/pq"
	"github.com/coopernurse/gorp"
	"database/sql"
	"log"
	"time"
	"./db/models"
	"strconv"
	"os"
)

func main() {
	db, err := sql.Open("postgres", "user=admin dbname=martinitodos sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	
	m := martini.Classic()
	m.Map(db)
	m.Use(martini.Static("assets"))
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.AddTableWithName(models.Todo{}, "todos").SetKeys(true, "Id")
	defer dbmap.Db.Close()

	m.Use(render.Renderer(render.Options {
		Layout: "layout",
	}))

	m.Get("/", func(r render.Render) {
		r.HTML(200, "welcome", map[string]interface{}{"greeting": "everybody"})
	})

	m.Get("/todos", func(r render.Render) {
		var todos []models.Todo
		_, err = dbmap.Select(&todos, "select * from todos order by created")
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(todos)
		r.HTML(200, "todos/index", todos)
	})

	m.Get("/todos/new", func(r render.Render) {
		r.HTML(200, "todos/new", nil)
	})

	m.Post("/todos", binding.Bind(models.Todo{}), func(todo models.Todo, r render.Render) {
		dbmap.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds)) 
		t1 := &models.Todo{Title: todo.Title, Description: todo.Description, Created: time.Now().UnixNano()}
		err := dbmap.Insert(t1)
		dbmap.TraceOff()
		if err != nil {
			log.Println(err)
		}
		r.Redirect("todos/" + strconv.FormatInt(t1.Id, 10), 302)
	})

	m.Get("/todos/:id", func(params martini.Params, r render.Render) {
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
