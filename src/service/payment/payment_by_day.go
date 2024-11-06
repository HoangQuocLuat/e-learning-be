package service_payment

import (
	"context"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	model_payment "e-learning/src/database/model/payment"
	"e-learning/src/service"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type PaymentGetByDayCommand struct {
	Month string
	Year  string
}

func PaymentGetByDay(ctx context.Context, s *PaymentGetByDayCommand) (total int, err error) {
	log.Println("[service_user.UserPagination] start")
	defer func() {
		log.Println("[service_user.UserPagination] end", "data", map[string]interface{}{"command: ": s}, "error", err)
	}()
	startDate, _ := time.Parse("2006-01-02", fmt.Sprintf("%s-%s-01", s.Year, s.Month)) // ngày đầu tiên của tháng
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)
	req := bson.M{
		"created_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
		"status": "Thanh toán thành công",
	}
	cursor, err := collection.Payment().Collection().Find(ctx, req)
	if err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Tuition + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return 0, fmt.Errorf(codeErr)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var payment model_payment.Payment
		err := cursor.Decode(&payment)
		if err != nil {
			codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Tuition + src_const.InternalError
			service.AddError(ctx, "", "", codeErr)
			return 0, fmt.Errorf(codeErr)
		}
		num, err := strconv.Atoi(payment.Amount)
		if err != nil {
			codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Tuition + src_const.InternalError
			service.AddError(ctx, "", "", codeErr)
			return 0, fmt.Errorf(codeErr)
		}
		total = total + num
	}
	return total, nil
}
