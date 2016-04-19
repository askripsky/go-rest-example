This is an example of a REST API implemented in Go.  Gin is used for the web framework and Ginkgo is used for tests.

## Referenced Posts
 * http://ewanvalentine.io/writing-and-running-go-apis-in-docker/
 * http://modocache.io/restful-go
 * https://blog.pivotal.io/labs/labs/a-rubyist-leaning-go-testing-http-handlers
 * https://blog.gopheracademy.com/advent-2013/day-11-martini/

## Running locally
go get github.com/askripsky/go-rest-example

cd $GOPATH/src/github.com/askripsky/go-rest-example
docker-compose build && docker-compose up

## Running tests

```
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega

go test
```

NOTE: If it doesn't build, make sure you have GOPATH setup correctly, i.e.
```
export GOPATH=$HOME
```