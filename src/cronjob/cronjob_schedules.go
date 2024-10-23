package cronjob

import (
	"context"
	"e-learning/src/database/collection"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/gomail.v2"
)

func NotifyWithTimeBySchedules() {
	c := cron.New()
	c.AddFunc("37 15 * * *", func() {
		log.Print("Cron job started")
		schedules()
	})

	c.Start()
	select {}
}

func schedules() {
	now := time.Now()
	dayOfWeek := int(now.Weekday())
	if dayOfWeek == 0 {
		dayOfWeek = 7 // Change Sunday to 7
	}
	log.Printf("Current time: %v", now)

	filter := bson.M{
		"day_of_week": dayOfWeek,
		"start_date": bson.M{
			"$lte": now,
		},
		"end_date": bson.M{
			"$gte": now,
		},
	}
	cursor, err := collection.Schedules().Collection().Find(context.TODO(), filter)
	if err != nil {
		log.Printf("Error finding schedules: %v", err)
		return // Return early on error
	}
	defer cursor.Close(context.TODO()) // Ensure cursor is closed

	log.Println("Cursor:", cursor)

	for cursor.Next(context.TODO()) {
		var schedule struct {
			ClassID     string    `bson:"class_id"`
			ClassName   string    `bson:"class_name"`
			Description string    `bson:"description"`
			StartTime   time.Time `bson:"start_time"`
			EndTime     time.Time `bson:"end_time"`
		}
		if err := cursor.Decode(&schedule); err != nil {
			log.Printf("Error decoding schedule: %v", err)
			continue
		}
		log.Printf("Schedule found: %+v", schedule)
		nowTime := time.Date(0, 1, 1, now.Hour(), now.Minute(), 0, 0, schedule.StartTime.Location())
		scheduleTime := time.Date(0, 1, 1, schedule.StartTime.Hour(), schedule.StartTime.Minute(), schedule.StartTime.Second(), schedule.StartTime.Nanosecond(), schedule.StartTime.Location())

		log.Println("Current time (hour-based):", nowTime)
		log.Println("Schedule start time (hour-based):", scheduleTime)

		if scheduleTime.Sub(nowTime) <= time.Hour && scheduleTime.Sub(nowTime) >= 0 {
			subject := fmt.Sprintf("Reminder: Class %s is starting soon", schedule.ClassName)
			body := fmt.Sprintf("Class %s will start at %s.", schedule.ClassName, schedule.StartTime.Format("15:04"))
			if err := SendMail("hoangquocluatspak@gmail.com", subject, body); err != nil {
				log.Println("Error sending email:", err)
			}
		}

		if schedule.EndTime.Sub(now) <= 30*time.Minute && schedule.EndTime.Sub(now) > 0 {
			subject := fmt.Sprintf("Reminder: Class %s is ending soon", schedule.ClassName)
			body := fmt.Sprintf("Class %s will end at %s.", schedule.ClassName, schedule.EndTime.Format("15:04"))
			if err := SendMail("hoangquocluatspak@gmail.com", subject, body); err != nil {
				log.Println("Error sending email:", err)
			}
		}
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
	}
}

func SendMail(email string, header string, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", "hoangquocluatspak@gmail.com")
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", header)
	msg.SetBody("text/plain", body)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, "hoangquocluatspak@gmail.com", "tyyk yafp tpdr qgio")

	if err := dialer.DialAndSend(msg); err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}
	return nil
}
