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
	newxsrf, newlses, newbear, err := lms.GetNewCookie(profile.Xsrf, profile.Lsession, config.Mongoconn)
	if err != nil {
		resp.Response = err.Error()
		helper.WriteResponse(respw, httpstatus, resp)
		return
	}
	resp.Info = newxsrf + "|" + newlses + "|" + newbear
	helper.WriteResponse(respw, httpstatus, resp)
}
