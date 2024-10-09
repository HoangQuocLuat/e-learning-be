package service_rest_zalo_payment

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	service_rest_req "e-learning/src/service/service.rest/request"
	service_rest_resp "e-learning/src/service/service.rest/response"

	"github.com/zpmep/hmacutil"
)

type object map[string]interface{}

var (
	app_id = "2553"
	key1   = "PcY4iZIKFCIdgZvA6ueMcMHHUbRLYjPL"
	// key2   = "kLtgPl8HHhfvMuDHPwKfgfsY4Ydm9eIz"
)

func Order(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rand.Seed(time.Now().UnixNano())
	transID := rand.Intn(1000000) // Generate random trans id
	embedData, _ := json.Marshal(object{})
	items, _ := json.Marshal([]object{})
	// request data
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req service_rest_req.Order
	if err := json.Unmarshal(body, &req); err != nil {
		log.Fatal(err)
	}

	params := make(url.Values)
	params.Add("app_id", app_id)
	params.Add("amount", req.Amount)
	params.Add("app_user", "user123")
	params.Add("embed_data", string(embedData))
	params.Add("item", string(items))
	params.Add("description", "Lazada - Payment for the order #"+strconv.Itoa(transID))
	params.Add("bank_code", "")

	now := time.Now()
	params.Add("app_time", strconv.FormatInt(now.UnixNano()/int64(time.Millisecond), 10)) // miliseconds

	params.Add("app_trans_id", fmt.Sprintf("%02d%02d%02d_%v", now.Year()%100, int(now.Month()), now.Day(), transID)) // translation missing: vi.docs.shared.sample_code.comments.app_trans_id

	// appid|app_trans_id|appuser|amount|apptime|embeddata|item
	data := fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v", params.Get("app_id"), params.Get("app_trans_id"), params.Get("app_user"),
		params.Get("amount"), params.Get("app_time"), params.Get("embed_data"), params.Get("item"))

	params.Add("mac", hmacutil.HexStringEncode(hmacutil.SHA256, key1, data))

	// Content-Type: application/x-www-form-urlencoded
	fmt.Println("params:", params)
	res, err := http.PostForm("https://sb-openapi.zalopay.vn/v2/create", params)

	// parse response
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body2, _ := io.ReadAll(res.Body)

	var result map[string]interface{}

	if err := json.Unmarshal(body2, &result); err != nil {
		log.Fatal(err)
	}
	responseData := service_rest_resp.Response{
		Status:  "success",
		Message: "Request processed successfully",
		Data:    result,
	}
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
