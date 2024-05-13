package helper

import (
	"strings"

	"github.com/gocroot/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func WebHook(WAKeyword, WAPhoneNumber, WAAPIQRLogin, WAAPIMessage string, msg model.IteungMessage, db *mongo.Database) (resp model.Response) {
	if IsLoginRequest(msg, WAKeyword) { //untuk whatsauth request login
		resp = HandlerQRLogin(msg, WAKeyword, WAPhoneNumber, db, WAAPIQRLogin)
	} else { //untuk membalas pesan masuk
		resp = HandlerIncomingMessage(msg, WAPhoneNumber, db, WAAPIMessage)
	}
	return
}

func RefreshToken(dt *model.WebHook, WAPhoneNumber, WAAPIGetToken string, db *mongo.Database) (res *mongo.UpdateResult, err error) {
	resp, err := PostStructWithToken[model.User]("Token", WAAPIToken(WAPhoneNumber, db), dt, WAAPIGetToken)
	if err != nil {
		return
	}
	profile := &model.Profile{
		Phonenumber: resp.PhoneNumber,
		Token:       resp.Token,
	}
	res, err = ReplaceOneDoc(db, "profile", bson.M{"phonenumber": resp.PhoneNumber}, profile)
	if err != nil {
		return
	}
	return
}

func IsLoginRequest(msg model.IteungMessage, keyword string) bool {
	return strings.Contains(msg.Message, keyword) && msg.From_link
}

func GetUUID(msg model.IteungMessage, keyword string) string {
	return strings.Replace(msg.Message, keyword, "", 1)
}

func HandlerQRLogin(msg model.IteungMessage, WAKeyword string, WAPhoneNumber string, db *mongo.Database, WAAPIQRLogin string) (resp model.Response) {
	dt := &model.WhatsauthRequest{
		Uuid:        GetUUID(msg, WAKeyword),
		Phonenumber: msg.Phone_number,
		Delay:       msg.From_link_delay,
	}
	resp, _ = PostStructWithToken[model.Response]("Token", WAAPIToken(WAPhoneNumber, db), dt, WAAPIQRLogin)
	return
}

func HandlerIncomingMessage(msg model.IteungMessage, WAPhoneNumber string, db *mongo.Database, WAAPIMessage string) (resp model.Response) {
	dt := &model.TextMessage{
		To:       msg.Chat_number,
		IsGroup:  false,
		Messages: GetRandomReplyFromMongo(msg, db),
	}
	if msg.Chat_server == "g.us" { //jika pesan datang dari group maka balas ke group
		dt.IsGroup = true
	}
	if (msg.Phone_number != "628112000279") && (msg.Phone_number != "6283131895000") { //ignore pesan datang dari iteung
		pasetoken := WAAPIToken(WAPhoneNumber, db)
		if pasetoken != "" {
			resp, _ = PostStructWithToken[model.Response]("Token", WAAPIToken(WAPhoneNumber, db), dt, WAAPIMessage)
		} else {
			resp.Response = "Not Found Phonenumber " + WAPhoneNumber + " in the profile collection db, check phone number env."
		}

	}
	return
}

func GetRandomReplyFromMongo(msg model.IteungMessage, db *mongo.Database) string {
	rply, _ := GetRandomDoc[model.Reply](db, "reply", 1)
	replymsg := strings.ReplaceAll(rply[0].Message, "#BOTNAME#", msg.Alias_name)
	replymsg = strings.ReplaceAll(replymsg, "\\n", "\n")
	return replymsg
}

func WAAPIToken(phonenumber string, db *mongo.Database) string {
	filter := bson.M{"phonenumber": phonenumber}
	apitoken, _ := GetOneDoc[model.Profile](db, "profile", filter)
	return apitoken.Token
}
