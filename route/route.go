package route

import (
	"net/http"

	"github.com/gocroot/controller"
)

func URL(w http.ResponseWriter, r *http.Request) {
	var method, path string = r.Method, r.URL.Path
	switch {
	case method == "GET" && path == "/":
		controller.GetHome(w, r)
	case method == "POST" && path == "/webhook/inbox":
		controller.PostInbox(w, r)
	case method == "GET" && path == "/refresh/token":
		controller.GetNewToken(w, r)
	default:
		controller.NotFound(w, r)
	}
}
