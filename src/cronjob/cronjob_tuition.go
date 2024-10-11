package cronjob

import (
	"context"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	model_user "e-learning/src/database/model/user"
	"e-learning/src/service"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func ComputeTuition() {
	c := cron.New()
	c.AddFunc("07 11 * * *", func() {
		log.Print("Cron job started")
		tuition()
	})

	c.Start()
	select {}
}

func tuition() {
	// bbbbbbbb
	ctx := context.Background()
	//thời gian tháng trước
	firstOfThisMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
	firstOfLastMonth := firstOfThisMonth.AddDate(0, -1, 0) // Lùi lại 1 tháng
	lastOfLastMonth := firstOfThisMonth.Add(-time.Second)  // lùi 1s lấy thời điểm tháng trước
	fmt.Println("aa", firstOfThisMonth, "bb", firstOfLastMonth, "cc", lastOfLastMonth)
	// lọc số buổi đi học trong tháng
	cursor, err := collection.User().Collection().Find(ctx, bson.M{})
	if err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_User + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user model_user.User
		err := cursor.Decode(&user)
		if err != nil {
			codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Class + src_const.InternalError
			service.AddError(ctx, "", "", codeErr)
			return
		}
		// results = append(results, &user)
		fmt.Println(user.ID)
		fil01 := bson.M{
			"user_id": user.ID,
			"time_check_in": bson.M{
				"$gte": firstOfLastMonth,
				"$lt":  lastOfLastMonth,
			},
			"status_check_in": "1",
		}
		c, err := collection.Attendance().Collection().CountDocuments(ctx, fil01)
		if err != nil {
			codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Class + src_const.InternalError
			service.AddError(ctx, "", "", codeErr)
			return
		}

		fmt.Println(c)

	}
	// đưa ra số buổi, tính tổng học phí, phí giảm và tiền cuối cùng
	// lưu thông tin vào database
	// gửi mail thông báo đóng tiền học phí
}
