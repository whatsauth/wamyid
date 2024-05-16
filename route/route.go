package route

import (
	"log"
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/controller"
	"github.com/gocroot/helper"
)

func URL(w http.ResponseWriter, r *http.Request) {
	if config.ErrorMongoconn != nil {
		log.Println(config.ErrorMongoconn.Error())
	}

	var method, path string = r.Method, r.URL.Path
	switch {
	case method == "GET" && path == "/":
		controller.GetHome(w, r)
	case method == "GET" && path == "/refresh/token":
		controller.GetNewToken(w, r)
	case method == "POST" && helper.URLParam(path, "/webhook/nomor/:nomorwa"):
		controller.PostInboxNomor(w, r)
	default:
		controller.NotFound(w, r)
	}
}
