package module

import (
	"github.com/gocroot/mod"
	"github.com/gocroot/model"
)

func Caller(Modulename string, Pesan model.IteungMessage) (reply string) {
	switch Modulename {
	case "idgrup":
		reply = mod.IDGroup(Pesan)
	}
	return
}
