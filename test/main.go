// go version go1.11.1 linux/amd64
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/zpmep/hmacutil" // go get github.com/zpmep/hmacutil
)

var (
	appID      = 2553
	key1       = "PcY4iZIKFCIdgZvA6ueMcMHHUbRLYjPL"
	appTransID = "241014_530740" // Input your app trans id
)

func main() {
	data := fmt.Sprintf("%v|%s|%s", appID, appTransID, key1) // appid|apptransid|key1
	fmt.Println("app", appID, "trans", appTransID, "key1", key1)
	params := map[string]interface{}{
		"app_id":       appID,
		"app_trans_id": appTransID,
		"mac":          hmacutil.HexStringEncode(hmacutil.SHA256, key1, data),
	}

	jsonStr, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.Post("https://sb-openapi.zalopay.vn/v2/query", "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var result map[string]interface{}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal(err)
	}

	for k, v := range result {
		log.Printf("%s = %+v", k, v)
	}
}
