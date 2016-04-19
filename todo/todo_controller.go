package todo

import (
    "labix.org/v2/mgo"
    "labix.org/v2/mgo/bson"
    "github.com/gin-gonic/gin"
)

type (
    TodoController struct {
        collection string
    }
)

func NewTodoController() *TodoController {
    return &TodoController{collection: "todos"}
}

func (controller *TodoController) GetAllTodos(c *gin.Context) {

    db := c.MustGet("db").(*mgo.Database)
    todos := []Todo{}

    err := db.C(controller.collection).Find(nil).Limit(100).All(&todos)

    if err != nil {
        panic(err)
    }

    c.JSON(200, todos)
}

func (controller *TodoController) InsertTodo(c *gin.Context) {

    var todo Todo
    c.BindJSON(&todo)

    db := c.MustGet("db").(*mgo.Database)

    todo.Id = bson.NewObjectId()
    db.C(controller.collection).Insert(todo)

    c.JSON(201, todo)
}

func (controller *TodoController) UpdateTodo(c *gin.Context) {

    var todo Todo
    c.BindJSON(&todo)

    db := c.MustGet("db").(*mgo.Database)
    id := c.Param("id")

    if !bson.IsObjectIdHex(id) {
        c.Status(404)
        return
    }

    err := db.C(controller.collection).UpdateId(bson.ObjectIdHex(id), todo)

    if err != nil {
        c.Status(404)
        return
    }

    c.JSON(200, todo)
}

func (controller *TodoController) DeleteTodo(c *gin.Context) {

    db := c.MustGet("db").(*mgo.Database)
    id := c.Param("id")

    if !bson.IsObjectIdHex(id) {
        c.Status(404)
        return
    }

    err := db.C(controller.collection).RemoveId(bson.ObjectIdHex(id))

    if err != nil {
        c.Status(404)
        return
    }

    c.Status(200)
}
