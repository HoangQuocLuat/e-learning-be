package service_order

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"net/url"
// 	"strconv"
// 	"time"

// 	src_const "e-learning/src/const"
// 	"e-learning/src/database/collection"
// 	model_order "e-learning/src/database/model/order"
// 	model_tuition "e-learning/src/database/model/tuition"
// 	"e-learning/src/service"

// 	"github.com/zpmep/hmacutil"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"golang.org/x/exp/rand"
// )

// type ItemInput struct {
// 	ItemID string
// 	Price  string
// }

// type OrderAddCommand struct {
// 	UserID string
// 	Amount string
// 	Items  []ItemInput
// }

// func OrderAdd(ctx context.Context, c *OrderAddCommand) (result *model_order.Order, err error) {
// 	//lấy ra tuition theo user_id
// 	tuition := &model_tuition.Tuition{}
// 	err = collection.Tuition().Collection().FindOne(ctx, bson.M{"user_id": c.UserID}).Decode(tuition)
// 	if err != nil {
// 		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Class + src_const.InternalError
// 		service.AddError(ctx, "", "", codeErr)
// 		return nil, fmt.Errorf(codeErr)
// 	}
// 	fmt.Println("aaa", tuition)
// 	// chuyển thông tin cho zalo pay
// 	rand.Seed(time.Now().UnixNano())
// 	transID := rand.Intn(1000000)

// 	params := make(url.Values)
// 	params.Add("app_id", src_const.App_id)
// 	params.Add("amount", string(tuition.Discount))
// 	params.Add("app_user", c.UserID)
// 	params.Add("embed_data", string(""))
// 	params.Add("item", string(""))
// 	params.Add("description", "Thanh toan hoc phi #"+strconv.Itoa(transID))
// 	params.Add("bank_code", "zalopayapp")
// 	now := time.Now()
// 	params.Add("app_time", strconv.FormatInt(now.UnixNano()/int64(time.Millisecond), 10))

// 	params.Add("app_trans_id", fmt.Sprintf("%02d%02d%02d_%v", now.Year()%100, int(now.Month()), now.Day(), transID))

// 	// appid|app_trans_id|appuser|amount|apptime|embeddata|item
// 	data := fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v", params.Get("app_id"), params.Get("app_trans_id"), params.Get("app_user"),
// 		params.Get("amount"), params.Get("app_time"), params.Get("embed_data"), params.Get("item"))

// 	params.Add("mac", hmacutil.HexStringEncode(hmacutil.SHA256, src_const.Key1, data))

// 	// Content-Type: application/x-www-form-urlencoded
// 	res, err := http.PostForm("https://sb-openapi.zalopay.vn/v2/create", params)

// 	// parse response
// 	var dataRes map[string]interface{}
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer res.Body.Close()
// 	r, _ := io.ReadAll(res.Body)
// 	if err := json.Unmarshal(r, &dataRes); err != nil {
// 		log.Fatal(err)
// 	}

// 	result = {
// 		OrderURL: dataRes.,
// 	}

// 	return nil, nil
// }
