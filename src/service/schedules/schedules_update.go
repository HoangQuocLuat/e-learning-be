package service_schedules

import (
	"context"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	model_schedules "e-learning/src/database/model/schedules"
	kafka_config "e-learning/src/kafka"
	"e-learning/src/service"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
)

type SchedulesUpdateCommand struct {
	ID string
	// DayOfWeek     *int
	// StartDate     *string
	// EndDate       *string
	StartTime     *string
	EndTime       *string
	Description   *string
	SchedulesType *int
}

func (t *SchedulesUpdateCommand) Valid() error {
	if t.ID == "" {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Tuition + src_const.InvalidErr
		return fmt.Errorf(codeErr)
	}
	return nil
}
func SchedulesUpdate(ctx context.Context, t *SchedulesUpdateCommand) (result *model_schedules.Schedules, err error) {
	if err := t.Valid(); err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Tuition + src_const.InvalidErr
		return nil, fmt.Errorf(codeErr)
	}
	err = collection.Schedules().Collection().FindOne(ctx, bson.M{"_id": t.ID}).Decode(&result)

	if err != nil {
		log.Println("[service_tuition.TuitionUpdate]", "FindOne ID", map[string]interface{}{"command: ": t}, "error", err)
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Tuition + src_const.TuitionExist
		return nil, fmt.Errorf(codeErr)
	}

	updateSchedules := bson.M{}

	// if t.DayOfWeek != nil {
	// 	updateSchedules["day_of_week"] = *t.DayOfWeek
	// 	result.DayOfWeek = *t.DayOfWeek
	// }
	// if t.StartDate != nil {
	// 	startDate, err := time.Parse("02-01-2006", *t.StartDate)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("invalid StartDate format: %v", err)
	// 	}
	// 	updateSchedules["start_date"] = startDate

	// 	result.StartDate = startDate
	// }

	// if t.EndDate != nil {
	// 	endDate, err := time.Parse("02-01-2006", *t.EndDate)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("invalid EndDate format: %v", err)
	// 	}
	// 	updateSchedules["end_date"] = endDate
	// 	result.EndDate = endDate
	// }

	if t.StartTime != nil {
		startTime, err := time.Parse("15:04", *t.StartTime)
		if err != nil {
			return nil, fmt.Errorf("invalid StartTime format: %v", err)
		}
		updateSchedules["start_time"] = startTime
		result.StartTime = startTime
	}

	if t.EndTime != nil {
		endTime, err := time.Parse("15:04", *t.EndTime)
		if err != nil {
			return nil, fmt.Errorf("invalid StartTime format: %v", err)
		}
		updateSchedules["start_time"] = endTime
		result.StartTime = endTime
	}

	if t.Description != nil {
		updateSchedules["description"] = *t.Description
		result.Description = *t.Description
	}

	if t.SchedulesType != nil {
		sType := *t.SchedulesType
		updateSchedules["schedules_type"] = src_const.MapSchedulesType[sType]
		result.SchedulesType = src_const.MapSchedulesType[sType]
	}

	_, err = collection.Schedules().Collection().UpdateOne(ctx, bson.M{"_id": t.ID}, bson.M{"$set": updateSchedules})

	if err != nil {
		log.Println("[service_order.TuitionUpdate]", "Update", map[string]interface{}{"command: ": t}, "error", err)
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Tuition + src_const.InternalError
		return nil, fmt.Errorf(codeErr)
	}

	err = collection.Schedules().Collection().FindOne(ctx, bson.M{"_id": t.ID}).Decode(&result)
	if err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf(codeErr)
	}

	mails, err := GetEmailsByClassID(ctx, result.ClassID)
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
	return
}
