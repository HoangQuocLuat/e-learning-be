package service_tuition

import (
	"context"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	model_tuition "e-learning/src/database/model/tuition"
	"e-learning/src/service"
	"fmt"
	"log"

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
	err = collection.Tuition().Collection().FindOne(ctx, bson.M{"user_id": c.UserID}).Decode(&results)
	if err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf(codeErr)
	}
	return
}
