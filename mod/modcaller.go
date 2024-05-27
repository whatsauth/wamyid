package mod

import (
	"github.com/gocroot/mod/idgrup"
	"github.com/gocroot/model"
)

func Caller(Modulename string, Pesan model.IteungMessage) (reply string) {
	switch Modulename {
	case "idgrup":
		reply = idgrup.IDGroup(Pesan)
	}
	return
}
