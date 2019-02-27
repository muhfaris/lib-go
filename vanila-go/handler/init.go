package handler

import "github.com/globalsign/mgo"

var (
	dbPool *mgo.Database
)

func Init(db *mgo.Database) {
	dbPool = db
}
