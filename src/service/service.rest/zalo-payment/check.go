package service_rest_zalo_payment

import (
	"bytes"
	src_const "e-learning/src/const"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/zpmep/hmacutil"
)

func Check(appTranID string) string {
	var result zalopayRess
	// var ctx context.Context
	appID, err := strconv.Atoi(src_const.App_id)
	if err != nil {
		log.Println("err", err)
		return ""
	}

	data := fmt.Sprintf("%v|%s|%s", appID, appTranID, src_const.Key1)
	params := map[string]interface{}{
		"app_id":       appID,
		"app_trans_id": appTranID,
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
	switch result.ReturnCode {
	case 1:
		return "success"
	case 2:
		return "failed"
	case 3:
		return "pending"
	default:
		return "pending"
	}

}
