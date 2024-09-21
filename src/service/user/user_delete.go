package service_user

import (
	"context"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	"e-learning/src/service"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type UserDeleteCommand struct {
	ID string
}

func (c *UserDeleteCommand) Valid() error {
	if c.ID == "" {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.InvalidErr
		return fmt.Errorf(codeErr)
	}

	return nil
}

func UserDelete(ctx context.Context, c *UserDeleteCommand) (err error) {
	log.Println("[service_user.UserDelete] start")
	defer func() {
		log.Println("[service_user.UserDelete] end", "data", map[string]interface{}{"user: ": c}, "error", err)
	}()

	if err = c.Valid(); err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.InvalidErr
		service.AddError(ctx, "", "", codeErr)
		return fmt.Errorf(codeErr)
	}
	_, err = collection.User().Collection().DeleteOne(ctx, bson.M{"_id": c.ID})
	if err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return fmt.Errorf(codeErr)
	}
	return
}
