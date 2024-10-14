package service_rest_zalo_payment

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	model_tuition "e-learning/src/database/model/tuition"
	service_rest_req "e-learning/src/service/service.rest/request"
	service_rest_resp "e-learning/src/service/service.rest/response"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/zpmep/hmacutil"
)

type object map[string]interface{}

func Order(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ctx context.Context
	var result map[string]interface{}
	rand.Seed(time.Now().UnixNano())
	transID := rand.Intn(1000000)
	//đọc dữ liệu từ req
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("req err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var req service_rest_req.Order
	if err := json.Unmarshal(body, &req); err != nil {
		log.Fatal(err)
	}

	//lấy ra tuition theo user_id
	tuition := &model_tuition.Tuition{}
	err = collection.Tuition().Collection().FindOne(ctx, bson.M{"user_id": req.AppUser}).Decode(tuition)
	if err != nil {
		log.Println("find err ", err)
		return
	}
	fmt.Println("aaa", tuition)

	//data fix cứng
	embedData, _ := json.Marshal(object{})
	fixedItems := []map[string]interface{}{
		{
			"itemid":       "knb",
			"itemname":     "kim nguyen bao",
			"itemprice":    198400,
			"itemquantity": 1,
		},
	}
	items, _ := json.Marshal(fixedItems)

	now := time.Now()
	trannnID := fmt.Sprintf("%02d%02d%02d_%v", now.Year()%100, int(now.Month()), now.Day(), transID)
	// log.Println(tuition.Discount)
	// log.Println("app_id", src_const.App_id, " app_user", req.AppUser, " item_id", string(items), "amount ", strconv.Itoa(tuition.Discount))
	// request data cho url zalo pay
	params := make(url.Values)
	params.Add("app_id", src_const.App_id)
	params.Add("amount", strconv.Itoa(tuition.Discount))
	params.Add("app_user", req.AppUser)
	params.Add("embed_data", string(embedData))
	params.Add("item", string(items))
	params.Add("description", "Thanh toan hoc phi #"+strconv.Itoa(transID))
	params.Add("bank_code", "zalopayapp")
	params.Add("app_time", strconv.FormatInt(now.UnixNano()/int64(time.Millisecond), 10)) // miliseconds

	params.Add("app_trans_id", trannnID) // translation missing: vi.docs.shared.sample_code.comments.app_trans_id

	// appid|app_trans_id|appuser|amount|apptime|embeddata|item
	data := fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v", params.Get("app_id"), params.Get("app_trans_id"), params.Get("app_user"),
		params.Get("amount"), params.Get("app_time"), params.Get("embed_data"), params.Get("item"))

	params.Add("mac", hmacutil.HexStringEncode(hmacutil.SHA256, src_const.Key1, data))

	// Content-Type: application/x-www-form-urlencoded
	res, err := http.PostForm("https://sb-openapi.zalopay.vn/v2/create", params)

	// parse response
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	ress, _ := io.ReadAll(res.Body)
	if err := json.Unmarshal(ress, &result); err != nil {
		log.Fatal(err)
	}
	responseData := service_rest_resp.Response{
		Status:  "success",
		Message: "Request processed successfully",
		Data: map[string]interface{}{
			"transaction_id": trannnID,
			"result":         result,
		},
	}
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
