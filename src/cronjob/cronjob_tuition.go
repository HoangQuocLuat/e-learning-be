package cronjob

import (
	"context"
	"e-learning/src/database/collection"
	model_tuition "e-learning/src/database/model/tuition"
	model_user "e-learning/src/database/model/user"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ComputeTuition() {
	c := cron.New()
	c.AddFunc("47 10 * * *", func() {
		log.Print("Cron job started")
		tuition()
	})

	c.Start()
	select {}
}

func tuition() {
	ctx := context.Background()
	var t *model_tuition.Tuition
	var discount int64
	// Thời gian tháng trước
	firstOfThisMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
	firstOfLastMonth := firstOfThisMonth.AddDate(0, -1, 0) // Lùi lại 1 tháng
	lastOfLastMonth := firstOfThisMonth.Add(-time.Second)  // Lùi 1 giây lấy thời điểm tháng trước

	// Lọc số buổi đi học trong tháng
	cursor, err := collection.User().Collection().Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user model_user.User
		err := cursor.Decode(&user)
		if err != nil {
			log.Printf("Error decoding user: %v", err)
			return
		}

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
			log.Printf("Error counting attendance: %v", err)
			return
		}

		if c == 0 {
			continue
		}

		totalFee := c * 30000
		//kiểm tra thuộc loại nào
		if user.UserType == "Loại 1" {
			discount = totalFee - (totalFee * 10 / 100)
		}
		if user.UserType == "Loại 2" {
			discount = totalFee - (totalFee * 20 / 100)
		}
		if user.UserType == "" {
			discount = totalFee
		}

		t = &model_tuition.Tuition{
			ID:        primitive.NewObjectID().Hex(),
			UserID:    user.ID,
			TotalFee:  int(totalFee),
			Discount:  int(discount),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Tính tổng học phí, phí giảm và tiền cuối cùng
		i, err := collection.Tuition().Collection().InsertOne(ctx, t)
		if err != nil {
			log.Printf("Error inserting tuition: %v", err)
			return
		}
		fmt.Println("Inserted a single document: ", i)
		subject := fmt.Sprintf("Noti: Đã có học phí tháng %d - năm %d", firstOfThisMonth.Month(), firstOfThisMonth.Year())
		body := fmt.Sprintf("Học phí tháng này:- Tổng tiền học %dVnĐ- Số tiền cần trả sau khi giảm %dVnĐ", totalFee, discount)
		if err := SendMail("hoangquocluatspak@gmail.com", subject, body); err != nil {
			log.Println("Error sending email:", err)
		}
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
	}
}
