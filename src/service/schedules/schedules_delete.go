package service_schedules

import (
	"context"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	"e-learning/src/service"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type SchedulesDeleteCommand struct {
	ID string
}

func (c *SchedulesDeleteCommand) Valid() error {
	if c.ID == "" {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.InvalidErr
		return fmt.Errorf(codeErr)
	}

	return nil
}

func SchedulesDelete(ctx context.Context, c *SchedulesDeleteCommand) (err error) {
	log.Println("[service_user.UserDelete] start")
	defer func() {
		log.Println("[service_user.UserDelete] end", "data", map[string]interface{}{"user: ": c}, "error", err)
	}()

	if err = c.Valid(); err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.InvalidErr
		service.AddError(ctx, "", "", codeErr)
		return fmt.Errorf(codeErr)
	}
	_, err = collection.Schedules().Collection().DeleteOne(ctx, bson.M{"_id": c.ID})
	if err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return fmt.Errorf(codeErr)
	}
	return
}
