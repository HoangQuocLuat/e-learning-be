package service_rest_zalo_payment

import (
	"context"
	"e-learning/config"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	model_payment "e-learning/src/database/model/payment"
	model_tuition "e-learning/src/database/model/tuition"
	service_rest_req "e-learning/src/service/service.rest/request"
	service_rest_resp "e-learning/src/service/service.rest/response"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/zpmep/hmacutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	if tuition.RemainingFee == 0 {
		responseData := service_rest_resp.Response{
			Status:  "PaymentCompleted",
			Message: "Hoc phi da thanh toan",
			Data:    nil,
		}
		if err := json.NewEncoder(w).Encode(responseData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	//data fix cứng
	embedData, _ := json.Marshal(object{
		"redirecturl": config.Get().ApiRedirect,
	})
	items, _ := json.Marshal([]object{})
	now := time.Now()
	trannnID := fmt.Sprintf("%02d%02d%02d_%v", now.Year()%100, int(now.Month()), now.Day(), transID)
	params := make(url.Values)
	params.Add("app_id", src_const.App_id)
	params.Add("amount", strconv.Itoa(tuition.Discount))
	params.Add("app_user", req.AppUser)
	params.Add("embed_data", string(embedData))
	params.Add("item", string(items))
	params.Add("description", "Thanh toan hoc phi #"+strconv.Itoa(transID))
	params.Add("bank_code", "zalopayapp")
	params.Add("app_time", strconv.FormatInt(now.UnixNano()/int64(time.Millisecond), 10)) // miliseconds
	params.Add("callback_url", config.Get().ApiCallBack)

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
	//tạo payment để trạng thái là đang thanh toán
	payment := &model_payment.Payment{
		ID:        primitive.NewObjectID().Hex(),
		UserID:    req.AppUser,
		TuitionID: tuition.ID,
		Status:    src_const.MapPaymentStatus[1],
		TransID:   trannnID,
		Amount:    strconv.Itoa(tuition.Discount),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = collection.Payment().Collection().InsertOne(ctx, payment)
	if err != nil {
		log.Println("err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
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
