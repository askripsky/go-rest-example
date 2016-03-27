FROM golang:latest

# Copy the local package files to the container’s workspace.
ADD . /go/src/github.com/askripsky/go-rest-example

# Install our dependencies
RUN go get github.com/go-martini/martini
RUN go get github.com/martini-contrib/binding
RUN go get github.com/martini-contrib/render
RUN go get labix.org/v2/mgo
RUN go get labix.org/v2/mgo/bson

# Install api binary globally within container
RUN go install github.com/askripsky/go-rest-example

# Set binary as entrypoint
ENTRYPOINT /go/bin/go-rest-example

# Expose default port (3000)
EXPOSE 3000
