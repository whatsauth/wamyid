package mod

import (
	"github.com/gocroot/mod/idgrup"
	"github.com/whatsauth/itmodel"
)

func Caller(Modulename string, Pesan itmodel.IteungMessage) (reply string) {
	switch Modulename {
	case "idgrup":
		reply = idgrup.IDGroup(Pesan)
	}
	return
}
