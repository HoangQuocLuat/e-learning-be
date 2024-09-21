package service_schedules

import (
	"context"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	model_schedules "e-learning/src/database/model/schedules"
	"e-learning/src/service"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type SchedulesByIDCommand struct {
	ClassID string
}

func (c *SchedulesByIDCommand) Valid() error {
	if c.ClassID == "" {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InvalidErr
		return fmt.Errorf(codeErr)
	}

	return nil
}

func SchedulesByClassID(ctx context.Context, c *SchedulesByIDCommand) (results []*model_schedules.Schedules, err error) {
	log.Println("[service_schedules.SchedulesByID] start")
	defer func() {
		log.Println("[service_schedules.SchedulesByID] end", "data", map[string]interface{}{"schedules: ": c}, "error", err)
	}()

	if err = c.Valid(); err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InvalidErr
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf(codeErr)
	}

	cursor, err := collection.Schedules().Collection().Find(ctx, bson.M{"class_id": c.ClassID})
	if err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf(codeErr)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var detailEntry model_schedules.Schedules
		if err := cursor.Decode(&detailEntry); err != nil {
			codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InvalidErr
			service.AddError(ctx, "", "", codeErr)
			return results, fmt.Errorf(codeErr)
		}
		results = append(results, &detailEntry)
	}
	if err := cursor.Err(); err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InvalidErr
		service.AddError(ctx, "", "", codeErr)
		return results, fmt.Errorf(codeErr)
	}
	return
}
