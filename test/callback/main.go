// go version go1.11.1 linux/amd64
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/zpmep/hmacutil" // go get github.com/zpmep/hmacutil
)

// App config
var (
	key2 = "kLtgPl8HHhfvMuDHPwKfgfsY4Ydm9eIz"
)

func main() {
	mux := http.DefaultServeMux
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var cbdata map[string]interface{}
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&cbdata)

		requestMac := cbdata["mac"].(string)
		dataStr := cbdata["data"].(string)
		mac := hmacutil.HexStringEncode(hmacutil.SHA256, key2, dataStr)
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
			log.Println("update order's status = success where app_trans_id =", dataJSON["app_trans_id"])
		}

		// thông báo kết quả cho ZaloPay server
		resultJSON, _ := json.Marshal(result)
		fmt.Fprintf(w, "%s", resultJSON)
	})

	log.Println("Server is listening at port :8888")
	http.ListenAndServe(":8888", mux)
}
