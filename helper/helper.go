package helper

import (
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func URLParam(reqpath string, url string) bool {
	urls := strings.Split(url, ":")
	prefix := reqpath[:strings.LastIndex(reqpath, "/")+1]
	return prefix == urls[0]
}

func GetParam(r *http.Request) string {
	return r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
}

func GetAddress() (ipport string, network string) {
	port := os.Getenv("PORT")
	network = "tcp4"
	if port == "" {
		port = ":8080"
	} else if port[0:1] != ":" {
		ip := os.Getenv("IP")
		if ip == "" {
			ipport = ":" + port
		} else {
			if strings.Contains(ip, ".") {
				ipport = ip + ":" + port
			} else {
				ipport = "[" + ip + "]" + ":" + port
				network = "tcp6"
			}
		}
	}
	return
}

func removeInvisibleChars(text string) string {
	// Create a regular expression to match invisible characters
	re := regexp.MustCompile(`\p{C}`)

	// Replace all matches with an empty string
	return re.ReplaceAllString(text, "")
}

func removeZeroWidthSpaces(text string) string {
	// Create a regular expression to match specific zero-width characters
	re := regexp.MustCompile(`[\u200B\u200C\u200D\uFEFF]`)

	// Replace all matches with an empty string
	return re.ReplaceAllString(text, "")
}

func NormalizeHiddenChar(text string) string {
	return removeZeroWidthSpaces(removeInvisibleChars(text))
}

func GetIPaddress() string {

	resp, err := http.Get("https://icanhazip.com/")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}
