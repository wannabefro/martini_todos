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
	"html/template"
	"strconv"
)

type Todo struct {
	Id					int64
	Created			int64
	Title				string `form:"Title" binding:"required"`
	Description	string `form:"Description" binding:"required`
}

func main() {
	m := martini.Classic()
	m.Use(martini.Static("assets"))
	dbmap := initDb()
	defer dbmap.Db.Close()

	m.Use(render.Renderer(render.Options {
		Layout: "layout",
		Funcs: []template.FuncMap {
			{
				"formatTime": func(args ...interface{}) string { 
					t1 := time.Unix(args[0].(int64), 0)
					return t1.Format("Mon Jan _2 2006")
				},
			},
		},
	}))

	m.Get("/", func(r render.Render) {
		r.HTML(200, "welcome", map[string]interface{}{"greeting": "everybody"})
	})

	m.Get("/todos", func(r render.Render) {
		var todos []Todo
		_, err := dbmap.Select(&todos, "select * from todos order by created")
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(todos)
		r.HTML(200, "todos/index", todos)
	})

	m.Get("/todos/new", func(r render.Render) {
		r.HTML(200, "todos/new", nil)
	})

	m.Post("/todos", binding.Form(Todo{}), func(todo Todo, errors binding.Errors, r render.Render) {
		if errors != nil {
			r.HTML(422, "todos/new", errors)
		} else {
			t1 := &Todo{Title: todo.Title, Description: todo.Description, Created: time.Now().Unix()}
			err := dbmap.Insert(t1)
			if err != nil {
				log.Println(err)
			}
			r.Redirect("todos/" + strconv.FormatInt(t1.Id, 10), 302)
		}
	})

	m.Get("/todos/:id", func(params martini.Params, r render.Render) {
		var todo Todo
		err := dbmap.SelectOne(&todo, "select * from todos where id = :id",
		map[string]interface{} {"id": params["id"]})
		if err != nil {
			log.Println(err)
		}
		r.HTML(200, "todos/show", todo)
	})


	m.Run()
}

func initDb() *gorp.DbMap {
	db, err := sql.Open("postgres", "user=admin dbname=martinitodos sslmode=disable")
	checkErr(err, "postgres.Open failed")

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
