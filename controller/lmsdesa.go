package controller

import (
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/mod/lms"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
)

func GetNewTokenLMSDesa(respw http.ResponseWriter, req *http.Request) {
	var resp itmodel.Response
	httpstatus := http.StatusServiceUnavailable
	profile, err := atdb.GetOneDoc[lms.LoginProfile](config.Mongoconn, "lmscreds", bson.M{})
	if err != nil {
		resp.Response = err.Error()
		helper.WriteResponse(respw, httpstatus, resp)
		return
	}
	profile.Xsrf, profile.Lsession, profile.Bearer, err = lms.GetNewCookie(profile.Xsrf, profile.Lsession, config.Mongoconn)
	if err != nil {
		resp.Response = err.Error()
		helper.WriteResponse(respw, httpstatus, resp)
		return
	}
	_, err = atdb.ReplaceOneDoc(config.Mongoconn, "lmscreds", bson.M{"user": profile.Username}, profile)
	if err != nil {
		resp.Response = err.Error()
		helper.WriteResponse(respw, httpstatus, resp)
		return
	}
	resp.Info = "ok"
	helper.WriteResponse(respw, http.StatusOK, resp)
}
