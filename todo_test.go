package main_test

import (
	. "github.com/askripsky/go-rest-example/server"
	. "github.com/askripsky/go-rest-example/database"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
    "net/http/httptest"
    "net/http"
    "encoding/json"
    "github.com/askripsky/go-rest-example/todo"
    "labix.org/v2/mgo"
    "labix.org/v2/mgo/bson"
    "bytes"
    "github.com/gin-gonic/gin"
    "strings"
)

var _ = Describe("Todos", func() {
    var dbName string
    var session *mgo.Session
    var server *gin.Engine
    var request *http.Request
    var recorder *httptest.ResponseRecorder

    BeforeEach(func() {
        dbName = "todos"
        session = Database()
        server = NewServer()

        recorder = httptest.NewRecorder()
    })

    AfterEach(func() {
        session.DB(dbName).DropDatabase()
    })

    Describe("GET /todos", func() {

        BeforeEach(func() {
            request, _ = http.NewRequest("GET", "/todos", nil)
        })

        Context("no Todos exist", func() {
            It("returns no todos", func() {
                server.ServeHTTP(recorder, request)
                Expect(recorder.Code).To(Equal(200))
                Expect(strings.TrimSpace(recorder.Body.String())).To(Equal("[]"))
            })
        })

        Context("when Todos exist", func() {

            BeforeEach(func() {
                collection := session.DB("todos").C("todos")
                collection.Insert(todo.Todo{Id: bson.NewObjectId(), Complete:false, Title:"Something"})
                collection.Insert(todo.Todo{Id: bson.NewObjectId(), Complete:true, Title:"Something Else"})
            })

            It("returns existing Todos", func() {
                server.ServeHTTP(recorder, request)
                Expect(recorder.Code).To(Equal(200))

                var todos []todo.Todo
                json.Unmarshal(recorder.Body.Bytes(), &todos)
                Expect(len(todos)).To(Equal(2))

                todo := todos[0]
                Expect(todo.Complete).To(Equal(false))
                Expect(todo.Title).To(Equal("Something"))
            })
        })
    })

    Describe("POST /todos", func() {

        BeforeEach(func() {
            var json = []byte(`{"title":"a title", "complete":false}`)
            request, _ = http.NewRequest("POST", "/todos", bytes.NewBuffer(json))
            request.Header.Set("Content-Type", "application/json")
        })

        Context("Create new Todo", func() {
            It("returns saved Todo", func() {
                server.ServeHTTP(recorder, request)
                Expect(recorder.Code).To(Equal(201))

                var savedTodo todo.Todo
                json.Unmarshal(recorder.Body.Bytes(), &savedTodo)
                Expect(savedTodo.Complete).To(Equal(false))
                Expect(savedTodo.Title).To(Equal("a title"))
            })
        })
    })

    Describe("PUT /todos/:id", func() {

        Context("no Todos exist", func() {

            BeforeEach(func() {
                var json = []byte(`{"title":"a title", "complete":false}`)
                request, _ = http.NewRequest("PUT", "/todos/123", bytes.NewBuffer(json))
                request.Header.Set("Content-Type", "application/json")
            })

            It("returns not found", func() {
                server.ServeHTTP(recorder, request)
                Expect(recorder.Code).To(Equal(404))
            })
        })

        Context("update existing Todo", func() {

            var collection *mgo.Collection
            var idToUpdate bson.ObjectId

            BeforeEach(func() {
                todoToUpdate := todo.Todo{Id: bson.NewObjectId(), Title:"Something"}
                todoToUpdate.Complete = false

                idToUpdate = todoToUpdate.Id

                collection = session.DB("todos").C("todos")
                collection.Insert(todoToUpdate)

                todoToUpdate.Complete = true
                todoJson, _ := json.Marshal(todoToUpdate)
                request, _ = http.NewRequest("PUT", "/todos/" + idToUpdate.Hex(), bytes.NewBuffer(todoJson))
                request.Header.Set("Content-Type", "application/json")
            })

            It("returns updated Todo", func() {
                server.ServeHTTP(recorder, request)
                Expect(recorder.Code).To(Equal(200))

                var updatedTodo todo.Todo
                json.Unmarshal(recorder.Body.Bytes(), &updatedTodo)
                Expect(updatedTodo.Complete).To(Equal(true))
            })

            It("updates the Todo in the database", func() {
                server.ServeHTTP(recorder, request)

                var updatedTodo todo.Todo
                collection.FindId(idToUpdate).One(&updatedTodo)
                Expect(updatedTodo.Title).To(Equal("Something"))
                Expect(updatedTodo.Complete).To(Equal(true))
            })
        })
    })

    Describe("DELETE /todos/:id", func() {

        Context("no Todos exist", func() {

            BeforeEach(func() {
                var json = []byte(`{"title":"a title", "complete":false}`)
                request, _ = http.NewRequest("DELETE", "/todos/123", bytes.NewBuffer(json))
                request.Header.Set("Content-Type", "application/json")
            })

            It("returns not found", func() {
                server.ServeHTTP(recorder, request)
                Expect(recorder.Code).To(Equal(404))
            })
        })

        Context("delete existing Todo", func() {

            var collection *mgo.Collection
            var idToDelete bson.ObjectId

            BeforeEach(func() {
                todoToDelete := todo.Todo{Id: bson.NewObjectId(), Title:"Something"}

                idToDelete = todoToDelete.Id

                collection = session.DB("todos").C("todos")
                collection.Insert(todoToDelete)

                request, _ = http.NewRequest("DELETE", "/todos/" + idToDelete.Hex(), nil)
            })

            It("returns success", func() {
                server.ServeHTTP(recorder, request)
                Expect(recorder.Code).To(Equal(200))

                Expect(len(recorder.Body.Bytes())).To(Equal(0))
            })

            It("deletes the Todo in the database", func() {
                server.ServeHTTP(recorder, request)

                count, _ := collection.FindId(idToDelete).Count()
                Expect(count).To(Equal(0))
            })
        })
    })
})
