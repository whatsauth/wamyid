package helpdesk

import (
	"strconv"
	"strings"

	"github.com/gocroot/helper/atapi"
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
		tim := helpdesk.(string)
		if strings.Contains(msg, tim) {
			team = tim
			return
		}
		helpdeskslist = append(helpdeskslist, tim)
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
		scpe := scp.(string)
		if strings.Contains(Pesan.Message, scpe) {
			scope = scpe
			return
		}
		scopeslist = append(scopeslist, scpe)
	}
	return
}

// mendapatkan scope helpdesk dari pesan
func GetOperatorFromScopeandTeam(scope, team string, db *mongo.Database) (operator Helpdesk, err error) {
	filter := bson.M{
		"scope": scope,
		"team":  team,
	}
	operator, err = atdb.GetOneLowestDoc[Helpdesk](db, "helpdesk", filter, "jumlahantrian")
	if err != nil {
		return
	}
	operator.JumlahAntrian += 1
	filter = bson.M{
		"scope":        scope,
		"team":         team,
		"phonenumbers": operator.Phonenumbers,
	}
	_, err = atdb.ReplaceOneDoc(db, "helpdesk", filter, operator)
	if err != nil {
		return
	}
	return
}

// handling key word
func StartHelpdesk(Pesan itmodel.IteungMessage, db *mongo.Database) (reply string) {
	namateam, helpdeskslist, err := GetNamaTeamFromPesan(Pesan, db)
	if err != nil {
		return err.Error()
	}
	//suruh pilih nama team kalo tidak ada
	if namateam == "" {
		reply = "Silakan memilih helpdesk yang anda tuju:\n"
		for i, helpdesk := range helpdeskslist {
			no := strconv.Itoa(i + 1)
			reply += no + ". " + helpdesk + "\n" + "wa.me/62895601060000?text=bantuan+operator+" + helpdesk + "\n"
		}
		return
	}
	//suruh pilih scope dari bantuan team
	scope, scopelist, err := GetScopeFromTeam(Pesan, namateam, db)
	if err != nil {
		return err.Error()
	}
	//pilih scope jika belum
	if scope == "" {
		reply = "Silakan memilih jenis bantuan yang anda butuhkan dari operator " + namateam + ":\n"
		for i, scope := range scopelist {
			no := strconv.Itoa(i + 1)
			reply += no + ". " + scope + "\n" + "wa.me/62895601060000?text=bantuan+operator+" + namateam + "+" + scope + "\n"
		}
		return
	}
	//menuliskan pertanyaan bantuan
	user := User{
		Scope:        scope,
		Team:         namateam,
		Name:         Pesan.Alias_name,
		Phonenumbers: Pesan.Phone_number,
	}
	_, err = atdb.InsertOneDoc(db, "helpdeskuser", user)
	if err != nil {
		return err.Error()
	}
	reply = "Silahkan kak " + Pesan.Alias_name + " mengetik pertanyaan atau bantuan yang ingin dijawab oleh operator: "

	return
}

// handling non key word
func PenugasanOperator(Profile itmodel.Profile, Pesan itmodel.IteungMessage, db *mongo.Database) (reply string, err error) {
	user, err := atdb.GetOneLatestDoc[User](db, "helpdeskuser", bson.M{"phonenumbers": Pesan.Phone_number})
	if err != nil {
		return
	}
	if !user.Terlayani {
		user.Masalah += "\n" + Pesan.Message
		if user.Operator.Name == "" || user.Operator.Phonenumbers == "" {
			var op Helpdesk
			op, err = GetOperatorFromScopeandTeam(user.Scope, user.Team, db)
			if err != nil {
				return
			}
			user.Operator = op
		}
		_, err = atdb.ReplaceOneDoc(db, "helpdeskuser", bson.M{"_id": user.ID}, user)
		if err != nil {
			return
		}

		msgstr := "User " + user.Name + " (" + user.Phonenumbers + ")\nMeminta tolong kakak " + user.Operator.Name + " untuk mencarikan solusi dari masalahnya:\n" + user.Masalah + "\nSilahkan langsung kontak di nomor wa.me/" + user.Phonenumbers
		msgstr += "\n\nJika sudah teratasi mohon inputkan solusi yang sudah di berikan ke user melalui link berikut:\nwa.me/62895601060000?text=solusi+dari+operator+helpdesk+:+"
		dt := &itmodel.TextMessage{
			To:       user.Operator.Phonenumbers,
			IsGroup:  false,
			Messages: msgstr,
		}
		go atapi.PostStructWithToken[itmodel.Response]("Token", Profile.Token, dt, Profile.URLAPIText)
		reply = "Kakak kami hubungkan dengan operator kami yang bernama *" + user.Operator.Name + "* di nomor wa.me/" + user.Operator.Phonenumbers + "\nMohon tunggu sebentar kami akan kontak kakak melalui nomor tersebut.\n_Terima kasih_"

	}
	return

}
