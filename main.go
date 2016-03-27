package main

import (
    "github.com/askripsky/go-rest-example/server"
)

func main() {
    server := server.NewServer()
    server.Run()
}
