package idname

import (
	"fmt"

	"github.com/whatsauth/itmodel"
)

func IDName(Pesan itmodel.IteungMessage) (reply string) {
	longitude := fmt.Sprintf("%f", Pesan.Longitude)
	latitude := fmt.Sprintf("%f", Pesan.Latitude)

	return "Hai.. hai.. kakak atas nama:\n" + Pesan.Alias_name + "\nLongitude: " + longitude + "\nLatitude: " + latitude + "\nberhasil absen\nmakasih"
}
