package service_schedules

import (
	"context"
	"fmt"
	"time"

	"e-learning/src/database/collection"
	kafka_config "e-learning/src/kafka"
	"e-learning/src/service"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	src_const "e-learning/src/const"
	model_schedules "e-learning/src/database/model/schedules"
	model_user "e-learning/src/database/model/user"
)

type SchedulesAddCommand struct {
	ClassID       string
	Description   string
	StartDate     string
	EndDate       string
	StartTime     string
	EndTime       string
	SchedulesType int
	DayOfWeek     int
}

func SchedulesAdd(ctx context.Context, c *SchedulesAddCommand) (result []*model_schedules.Schedules, err error) {
	startDate, err := time.Parse("02-01-2006", c.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid StartDate format: %v", err)
	}

	endDate, err := time.Parse("02-01-2006", c.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid EndDate format: %v", err)
	}

	startTime, err := time.Parse("15:04", c.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid StartTime format: %v", err)
	}

	endTime, err := time.Parse("15:04", c.EndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid EndTime format: %v", err)
	}

	targetDay, err := IntToWeekday(c.DayOfWeek)
	if err != nil {
		return nil, err
	}

	for date := startDate; date.Before(endDate) || date.Equal(endDate); date = date.AddDate(0, 0, 1) {
		if date.Weekday() == targetDay {
			result = append(result, &model_schedules.Schedules{
				ID:            primitive.NewObjectID().Hex(),
				ClassID:       c.ClassID,
				Description:   c.Description,
				SchedulesType: src_const.MapSchedulesType[c.SchedulesType],
				DayOfWeek:     c.DayOfWeek,
				Day:           date,
				StartTime:     startTime,
				EndTime:       endTime,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			})
		}
	}
	var scheduleDocs []interface{}
	for _, schedule := range result {
		scheduleDocs = append(scheduleDocs, schedule)
	}

	_, err = collection.Schedules().Collection().InsertMany(ctx, scheduleDocs)
	if err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf(codeErr)
	}

	mails, err := GetEmailsByClassID(ctx, c.ClassID)
	if err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf(codeErr)
	}

	if mails != nil {
		var messages []kafka.Message
		for _, email := range mails {
			message := kafka.Message{
				Key:   []byte("send-email"),
				Value: []byte(email),
				Time:  time.Now(),
			}
			messages = append(messages, message)
		}
		// Ghi các tin nhắn vào Kafka
		err := kafka_config.KafkaProducer.WriteMessages(ctx, messages...)
		if err != nil {
			fmt.Println(err)
			codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InternalError
			service.AddError(ctx, "", "", codeErr)
			return nil, fmt.Errorf(codeErr)
		}
	}
	return result, nil
}

func GetEmailsByClassID(ctx context.Context, classID string) ([]string, error) {
	cursor, err := collection.User().Collection().Find(ctx, bson.M{"class_id": classID})
	if err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf("%s: %v", codeErr, err)
	}
	defer cursor.Close(ctx)

	var emails []string
	for cursor.Next(ctx) {
		var user model_user.User
		if err := cursor.Decode(&user); err != nil {
			codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InternalError
			service.AddError(ctx, "", "", codeErr)
			return nil, fmt.Errorf("%s: %v", codeErr, err)
		}
		emails = append(emails, user.Email)
	}

	if err := cursor.Err(); err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf("%s: %v", codeErr, err)
	}
	return emails, nil
}

func IntToWeekday(day int) (time.Weekday, error) {
	if day < 1 || day > 7 {
		return time.Sunday, fmt.Errorf("invalid day number: %d (must be 1-7)", day)
	}
	return time.Weekday(day % 7), nil
}
