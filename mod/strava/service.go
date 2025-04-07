package strava

import (
	"errors"
	"net/http"

	"github.com/gocroot/helper/atapi"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/mod/pomokit"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getConfigByPhone(db *mongo.Database, profilePhone string) (*Config, error) {
	conf, err := atdb.GetOneDoc[Config](db, "config", bson.M{"phonenumber": profilePhone})
	if err != nil {
		return nil, errors.New("kesalahan dalam pengambilan config di database: " + err.Error())
	}
	return &conf, nil
}

func postToDomyikado(secret, url string, data map[string]interface{}) error {
	statuscode, httpresp, err := atapi.PostStructWithToken[itmodel.Response]("secret", secret, data, url)
	if err != nil {
		return errors.New("akses ke endpoint domyikado gagal: " + err.Error())
	}

	if statuscode != http.StatusOK {
		return errors.New("salah posting endpoint domyikado: " + httpresp.Response + "\ninfo: " + httpresp.Info)
	}

	return nil
}

func getWaGroupIDFromPomokit(db *mongo.Database, phoneNumber string) (string, error) {
	pomokitDoc, err := atdb.GetOneDoc[pomokit.PomodoroReport](db, "pomokit", bson.M{"phonenumber": phoneNumber})
	if err != nil {
		return "", errors.New("kesalahan dalam pengambilan data pomokit: " + err.Error())
	}

	if pomokitDoc.WaGroupID == "" {
		return "", errors.New("grup ID tidak ditemukan")
	}

	return pomokitDoc.WaGroupID, nil
}
