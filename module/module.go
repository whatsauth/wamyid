package module

import (
	"fmt"
	"strings"

	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetModuleName(im model.IteungMessage, MongoConn *mongo.Database, ModuleCollection string) (modulename string) {
	modules, _ := atdb.GetAllDoc[[]Module](MongoConn, ModuleCollection)
	for _, mod := range modules {
		complete, _ := IsMatch(strings.ToLower(im.Message), mod.Keyword...)
		if complete {
			modulename = mod.Name
		}
	}
	return
}

func IsMatch(str string, subs ...string) (bool, int) {

	matches := 0
	isCompleteMatch := true

	fmt.Printf("String: \"%s\", Substrings: %s\n", str, subs)

	for _, sub := range subs {
		if strings.Contains(str, sub) {
			matches += 1
		} else {
			isCompleteMatch = false
		}
	}

	return isCompleteMatch, matches
}
