package database

import (
	"labix.org/v2/mgo"
    "os"
)

func Database() *mgo.Session {
    server, isServerSet := os.LookupEnv("DB_SERVER")

    if (!isServerSet) {
        server = "localhost"
    }

	session, err := mgo.Dial(server)

	if err != nil {
		panic(err)
	}

	return session
}
