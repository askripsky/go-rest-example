package server

import (
	"github.com/askripsky/go-rest-example/database"
	"github.com/askripsky/go-rest-example/todo"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
    "os"
)

func NewServer() *martini.ClassicMartini {
	m := martini.Classic()
	m.Use(render.Renderer())
    m.Use(DbHandler())

	todoController := todo.NewTodoController()

	m.Get("/todos", binding.Bind(todo.Todo{}), todoController.GetAllTodos)
	m.Post("/todos", binding.Bind(todo.Todo{}), todoController.InsertTodo)
	m.Put("/todos/:id", binding.Bind(todo.Todo{}), todoController.UpdateTodo)
	m.Delete("/todos/:id", todoController.DeleteTodo)

    return m
}

// Clone Session and close for each request
func DbHandler() martini.Handler {
    session := database.Database();

    dbName, isNameSet := os.LookupEnv("DB_NAME")

    if (!isNameSet) {
        dbName = "todos"
    }

    return func(c martini.Context) {
        s := session.Clone()
        c.Map(s.DB(dbName))
        defer s.Close()
        c.Next()
    }
}
