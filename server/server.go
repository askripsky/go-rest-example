package server

import (
    "github.com/askripsky/go-rest-example/database"
    "github.com/askripsky/go-rest-example/todo"
    "github.com/gin-gonic/gin"
    "os"
)

func NewServer() *gin.Engine {
    s := gin.Default()
    s.Use(DbHandler())

    todoController := todo.NewTodoController()

    s.GET("/todos", todoController.GetAllTodos)
    s.POST("/todos", todoController.InsertTodo)
    s.PUT("/todos/:id", todoController.UpdateTodo)
    s.DELETE("/todos/:id", todoController.DeleteTodo)

    return s
}

// Clone Session and close for each request
func DbHandler() gin.HandlerFunc {
    session := database.Database();

    dbName, isNameSet := os.LookupEnv("DB_NAME")

    if (!isNameSet) {
        dbName = "todos"
    }

    return func(c *gin.Context) {
        s := session.Clone()
        c.Set("db", s.DB(dbName))
        defer s.Close()
        c.Next()
    }
}
