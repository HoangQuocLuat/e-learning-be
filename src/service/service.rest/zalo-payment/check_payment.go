package service_rest_zalo_payment

import (
	"bytes"
	"context"
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
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/zpmep/hmacutil"
)

type zalopayRess struct {
	Amount     int    `json:"amount"`
	ReturnCode int    `json:"return_code"`
	ZPTransID  int    `json:"zp_trans_id"`
	ReturnMess string `json:"return_message"`
}

func CheckPay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var result zalopayRess
	var ctx context.Context
	//đọc dữ liệu req
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
	appID, err := strconv.Atoi(src_const.App_id)
	if err != nil {
		log.Println("err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := fmt.Sprintf("%v|%s|%s", appID, req.AppTransID, src_const.Key1)
	params := map[string]interface{}{
		"app_id":       appID,
		"app_trans_id": req.AppTransID,
		"mac":          hmacutil.HexStringEncode(hmacutil.SHA256, src_const.Key1, data),
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

	ress, _ := io.ReadAll(res.Body)
	if err := json.Unmarshal(ress, &result); err != nil {
		log.Fatal(err)
	}

	// lưu vào data base thông tin giao dịch
	if result.ReturnCode == 1 {
		payment := &model_payment.Payment{
			ID:        primitive.NewObjectID().Hex(),
			UserID:    req.UserID,
			TuitionID: req.TuitionID,
			TransID:   strconv.Itoa(result.ZPTransID),
			Amount:    strconv.Itoa(result.Amount),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_, err = collection.Payment().Collection().InsertOne(ctx, payment)
		if err != nil {
			log.Println("err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// fil tuition và cập nhật lại tiền đã thanh toán và tiền nợ ( nên cho vào kafka )
		// đoạn code này cần được xem lại vì trừ lần đầu ok còn lần 2 có vẻ còn lỗi
		fil := bson.M{
			"_id": req.TuitionID,
		}
		var tuition model_tuition.Tuition
		err = collection.Tuition().Collection().FindOne(ctx, fil).Decode(&tuition)
		fmt.Println(tuition)
		if err != nil {
			log.Println("Error decoding user:", err)
			return
		}
		remainingFee := tuition.Discount - result.Amount
		update := bson.M{
			"$set": bson.M{
				"paid_amount":   result.Amount,
				"remaining_fee": remainingFee,
				"updated_at":    time.Now(),
			},
		}
		_, err = collection.Tuition().Collection().UpdateOne(ctx, fil, update)
		if err != nil {
			log.Println("Error updating user:", err)
			return
		}
		//
	}

	responseData := service_rest_resp.Response{
		Status:  "success",
		Message: "Request processed successfully",
		Data: map[string]interface{}{
			"status":      result.ReturnCode,
			"amount":      result.Amount,
			"trans_id":    result.ZPTransID,
			"return_mess": result.ReturnMess,
		},
	}

	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
