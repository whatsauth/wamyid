package config

import (
	"os"

	"github.com/gocroot/helper"
	"github.com/gocroot/model"
)

var MongoString string = os.Getenv("MONGOSTRING")

var mongoinfo = model.DBInfo{
	DBString: MongoString,
	DBName:   "iteung",
}

var Mongoconn, _ = helper.MongoConnect(mongoinfo)
