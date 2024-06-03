package mod

import (
	"github.com/gocroot/mod/idgrup"
	"github.com/gocroot/mod/idname"
	"github.com/whatsauth/itmodel"
)

func Caller(Modulename string, Pesan itmodel.IteungMessage) (reply string) {
	switch Modulename {
	case "idgrup":
		reply = idgrup.IDGroup(Pesan)
	case "idname":
		reply = idname.IDName(Pesan)
	}
	return
}
