package service_rest_zalo_payment

import (
	"context"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	model_payment "e-learning/src/database/model/payment"
	model_tuition "e-learning/src/database/model/tuition"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/zpmep/hmacutil"
	"go.mongodb.org/mongo-driver/bson"
)

func CallbackPayment(w http.ResponseWriter, r *http.Request) {
	var ctx context.Context
	defer r.Body.Close()
	var cbdata map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&cbdata)

	requestMac := cbdata["mac"].(string)
	log.Println("req_mac", requestMac)
	dataStr := cbdata["data"].(string)
	mac := hmacutil.HexStringEncode(hmacutil.SHA256, src_const.Key2, dataStr)
	log.Println("mac =", mac)

	result := make(map[string]interface{})

	// kiểm tra callback hợp lệ (đến từ ZaloPay server)
	if mac != requestMac {
		// callback không hợp lệ
		result["return_code"] = -1
		result["return_message"] = "mac not equal"
	} else {
		// thanh toán thành công
		result["return_code"] = 1
		result["return_message"] = "success"

		// merchant cập nhật trạng thái cho đơn hàng
		var dataJSON map[string]interface{}
		json.Unmarshal([]byte(dataStr), &dataJSON)

		//update database payment
		fil := bson.M{
			"trans_id": dataJSON["app_trans_id"],
		}
		update := bson.M{
			"$set": bson.M{
				"status":     src_const.MapPaymentStatus[0],
				"updated_at": time.Now(),
			},
		}
		_, err := collection.Payment().Collection().UpdateOne(ctx, fil, update)
		if err != nil {
			log.Println("Error updating user:", err)
			return
		}
		//update database tuition
		var payment model_payment.Payment
		err = collection.Payment().Collection().FindOne(ctx, bson.M{"trans_id": dataJSON["app_trans_id"]}).Decode(&payment)
		if err != nil {
			log.Println("Error decoding user:", err)
			return
		}
		var tuition model_tuition.Tuition
		err = collection.Tuition().Collection().FindOne(ctx, bson.M{
			"_id": payment.TuitionID,
		}).Decode(&tuition)
		if err != nil {
			log.Println("Error decoding user:", err)
			return
		}
		// convert
		var amount int

		// Chuyển đổi dataJSON["amount"] sang int
		if val, ok := dataJSON["amount"].(float64); ok {
			amount = int(val) // Chuyển đổi từ float64 sang int
		} else if val, ok := dataJSON["amount"].(int); ok {
			amount = val // Nếu đã là int thì không cần chuyển đổi
		} else {
			log.Println("Error: amount is not a float64 or int")
			return
		}
		remainingFee := tuition.Discount - amount
		updateTuition := bson.M{
			"$set": bson.M{
				"paid_amount":   amount,
				"remaining_fee": remainingFee,
				"updated_at":    time.Now(),
			},
		}
		_, err = collection.Tuition().Collection().UpdateOne(ctx, bson.M{
			"_id": payment.TuitionID,
		}, updateTuition)
		if err != nil {
			log.Println("Error updating user:", err)
			return
		}
		log.Println("update order's status = success where app_trans_id =", dataJSON["app_trans_id"])
	}

	// thông báo kết quả cho ZaloPay server
	resultJSON, _ := json.Marshal(result)
	fmt.Fprintf(w, "%s", resultJSON)
}
