package config

import "os"

var PublicKey string = os.Getenv("PUBLICKEY")
