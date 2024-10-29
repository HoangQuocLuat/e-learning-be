package service_tuition

import (
	"context"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	model_tuition "e-learning/src/database/model/tuition"
	"e-learning/src/service"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type TuitionByIDCommand struct {
	UserID string
}

func (c *TuitionByIDCommand) Valid() error {
	if c.UserID == "" {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InvalidErr
		return fmt.Errorf(codeErr)
	}

	return nil
}

func TuitionGetByID(ctx context.Context, c *TuitionByIDCommand) (results *model_tuition.Tuition, err error) {
	log.Println("[service_schedules.SchedulesByID] start")
	defer func() {
		log.Println("[service_schedules.SchedulesByID] end", "data", map[string]interface{}{"schedules: ": c}, "error", err)
	}()

	if err = c.Valid(); err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InvalidErr
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf(codeErr)
	}
	now := time.Now()

	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)
	find := bson.M{
		"user_id": c.UserID,
		"created_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}
	err = collection.Tuition().Collection().FindOne(ctx, find).Decode(&results)
	if err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf(codeErr)
	}
	return
}
