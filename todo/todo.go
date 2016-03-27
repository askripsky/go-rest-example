package todo

import (
	"labix.org/v2/mgo/bson"
)

type Todo struct {
	Id       bson.ObjectId `json:"id" bson:"_id"`
	Title    string        `form:"title" json:"title"`
	Complete bool          `form:"complete" json:"complete"`
}
