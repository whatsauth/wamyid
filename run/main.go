package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/whatsauth/itmodel"
)

func PanduanDosen(message itmodel.IteungMessage) string {
	// Path file panduan_dosen.txt
	const filePath = "../mod/siakad/panduan_dosen.txt"

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return "Maaf, terjadi kesalahan saat mengambil panduan dosen."
	}
	return string(content)
}

func main() {
	// Example message
	message := itmodel.IteungMessage{
		Message:      "minta panduan dosen",
		Alias_name:   "Dosen1",
		Phone_number: "6281234567890",
	}

	// Call PanduanDosen function
	response := PanduanDosen(message)
	fmt.Println(response)
}
