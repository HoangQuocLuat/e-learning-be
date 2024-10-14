package service_rest_zalo_payment

import (
	"bytes"
	src_const "e-learning/src/const"
	service_rest_req "e-learning/src/service/service.rest/request"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/zpmep/hmacutil"
)

func CheckPay(w http.ResponseWriter, r *http.Request) {
	//đọc dữ liệu req
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("req err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var req service_rest_req.CheckPay
	if err := json.Unmarshal(body, &req); err != nil {
		log.Fatal(err)
	}

	data := fmt.Sprintf("%v|%s|%s", src_const.App_id, req.AppTransID, src_const.Key1)
	params := map[string]interface{}{
		"app_id":       src_const.App_id,
		"app_trans_id": req.AppTransID,
		"mac":          hmacutil.HexStringEncode(hmacutil.SHA256, src_const.Key1, data),
	}
	fmt.Println(params, data, "", src_const.Key1)
	jsonStr, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.Post("https://sb-openapi.zalopay.vn/v2/query", "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body1, _ := ioutil.ReadAll(res.Body)

	var result map[string]interface{}

	if err := json.Unmarshal(body1, &result); err != nil {
		log.Fatal(err)
	}

	for k, v := range result {
		log.Printf("%s = %+v", k, v)
	}
}
