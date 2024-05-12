package config

import (
	"gocroot/helper"
	"gocroot/model"
	"os"
)

var MongoString string = os.Getenv("MONGOSTRING")

var mongoinfo = model.DBInfo{
	DBString: MongoString,
	DBName:   "iteung",
}

var Mongoconn, _ = helper.MongoConnect(mongoinfo)
