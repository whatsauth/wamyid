package helper

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func PostStructWithToken[T any](tokenkey string, tokenvalue string, structname interface{}, urltarget string) (result T, err error) {
	client := http.Client{}
	mJson, _ := json.Marshal(structname)
	req, err := http.NewRequest("POST", urltarget, bytes.NewBuffer(mJson))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(tokenkey, tokenvalue)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if err = json.Unmarshal(respBody, &result); err != nil {
		return
	}
	return
}
