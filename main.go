package gocroot

import (
	"gocroot/controller"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("WebHook", controller.HandleRequest)
}
