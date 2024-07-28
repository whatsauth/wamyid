package helpdesk

import (
	"strconv"
	"strings"

	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// mendapatkan nama team helpdesk dari pesan
func GetNamaTeamFromPesan(Pesan itmodel.IteungMessage, db *mongo.Database) (team string, helpdeskslist []string, err error) {
	msg := strings.ReplaceAll(Pesan.Message, "bantuan", "")
	msg = strings.ReplaceAll(msg, "operator", "")
	msg = strings.TrimSpace(msg)
	helpdesks, err := atdb.GetAllDistinctDoc(db, bson.M{}, "team", "helpdesk")
	if err != nil {
		return
	}
	//mendapatkan keyword masuk ke team yang mana
	for _, helpdesk := range helpdesks {
		team = helpdesk.(string)
		if strings.Contains(msg, team) {
			return
		}
		helpdeskslist = append(helpdeskslist, team)
	}
	return
}

// mendapatkan scope helpdesk dari pesan
func GetScopeFromTeam(Pesan itmodel.IteungMessage, team string, db *mongo.Database) (scope string, scopeslist []string, err error) {
	filter := bson.M{
		"team": team,
	}
	scopes, err := atdb.GetAllDistinctDoc(db, filter, "scope", "helpdesk")
	if err != nil {
		return
	}
	//mendapatkan keyword masuk ke team yang mana
	for _, scp := range scopes {
		scope = scp.(string)
		if strings.Contains(Pesan.Message, scope) {
			return
		}
		scopeslist = append(scopeslist, scope)
	}
	return
}

// mendapatkan scope helpdesk dari pesan
func GetOperatorFromScopeandTeam(Pesan itmodel.IteungMessage, scope, team string, db *mongo.Database) (operator Helpdesk, err error) {
	filter := bson.M{
		"scope":    scope,
		"team":     team,
		"occupied": false,
	}
	operator, err = atdb.GetOneDoc[Helpdesk](db, "helpdesk", filter)
	if err != nil {
		return
	}
	operator.Occupied = true
	filter = bson.M{
		"scope":        scope,
		"team":         team,
		"phonenumbers": operator.Phonenumbers,
		"occupied":     false,
	}
	_, err = atdb.ReplaceOneDoc(db, "helpdesk", filter, operator)
	if err != nil {
		return
	}
	return
}

func PenugasanOperator(Pesan itmodel.IteungMessage, db *mongo.Database) (reply string) {
	namateam, helpdeskslist, err := GetNamaTeamFromPesan(Pesan, db)
	if err != nil {
		return err.Error()
	}
	//suruh pilih nama team kalo tidak ada
	if namateam == "" {
		reply = "Silakan memilih helpdesk yang anda tuju:\n"
		for i, helpdesk := range helpdeskslist {
			no := strconv.Itoa(i + 1)
			reply += no + ". " + helpdesk + "\n" + "https://wa.me/62895601060000?text=bantuan+operator+" + helpdesk + "\n"
		}
		return
	}
	//suruh pilih scope dari bantuan team
	scope, scopelist, err := GetScopeFromTeam(Pesan, namateam, db)
	//pilih scope jika belum
	if scope == "" {
		reply = "Silakan memilih jenis bantuan yang anda butuhkan:\n"
		for i, scope := range scopelist {
			no := strconv.Itoa(i + 1)
			reply += no + ". " + scope + "\n" + "https://wa.me/62895601060000?text=bantuan+operator+" + scope + "\n"
		}
		return
	}
	//menuliskan pertanyaan bantuan
	user := Helpdesk{
		Scope:        scope,
		Team:         namateam,
		Name:         Pesan.Alias_name,
		Phonenumbers: Pesan.Phone_number,
		Occupied:     false, //false kalo belum dilayani
	}
	_, err = atdb.InsertOneDoc(db, "helpdeskuser", user)
	if err != nil {
		return err.Error()
	}
	reply = "Silahkan kak " + Pesan.Alias_name + " mengetik pertanyaan atau bantuan yang ingin dijawab oleh operator: "

	return
}
