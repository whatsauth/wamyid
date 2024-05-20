package gocroot

import (
	"testing"

	"github.com/gocroot/config"
	"github.com/gocroot/helper"
)

func TestIsBotNumber(t *testing.T) {
	result, err := helper.IsBotNumber("34324234", config.Mongoconn)
	print(result)
	print(err)
	result, err = helper.IsBotNumber("234324", config.Mongoconn)
	print(result)
	print(err)
}
