package kimseok

import (
	"testing"
)

// Fungsi ini akan dijalankan oleh `go test` dan memeriksa fungsi atau metode yang ingin Anda uji.
func TestExampleFunction(t *testing.T) {

	queries := Stemmer("cara ngoding golang ")
	print(queries)
}
