package todo

import (
	"github.com/martini-contrib/render"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
    "github.com/go-martini/martini"
)

type (
	TodoController struct {
        collection string
	}
)

func NewTodoController() *TodoController {
	return &TodoController{collection: "todos"}
}

func (controller *TodoController) GetAllTodos(r render.Render, db *mgo.Database) {
	todos := []Todo{}

	err := db.C(controller.collection).Find(nil).Limit(100).All(&todos)

	if err != nil {
		panic(err)
	}

	r.JSON(200, todos)
}

func (controller *TodoController) InsertTodo(todo Todo, r render.Render, db *mgo.Database) {

	todo.Id = bson.NewObjectId()
    db.C(controller.collection).Insert(todo)

	r.JSON(201, todo)
}

func (controller *TodoController) UpdateTodo(params martini.Params, todo Todo, r render.Render, db *mgo.Database) {

    id := params["id"]

    if !bson.IsObjectIdHex(id) {
        r.Status(404)
        return
    }

    err := db.C(controller.collection).UpdateId(bson.ObjectIdHex(id), todo)

    if err != nil {
        r.Status(404)
        return
    }

    r.JSON(200, todo)
}

func (controller *TodoController) DeleteTodo(params martini.Params, r render.Render, db *mgo.Database) {

    id := params["id"]

    if !bson.IsObjectIdHex(id) {
        r.Status(404)
        return
    }

    err := db.C(controller.collection).RemoveId(bson.ObjectIdHex(id))

    if err != nil {
        r.Status(404)
        return
    }

    r.Status(200)
}
