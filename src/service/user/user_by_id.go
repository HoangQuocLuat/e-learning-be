package service_user

import (
	"context"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	model_user "e-learning/src/database/model/user"
	"e-learning/src/service"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type UserByIDCommand struct {
	UserID string
}

func (c *UserByIDCommand) Valid() error {
	if c.UserID == "" {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.InvalidErr
		return fmt.Errorf(codeErr)
	}

	return nil
}

func UserByID(ctx context.Context, c *UserByIDCommand) (result *model_user.User, err error) {
	log.Println("[service_user.UserByID] start")
	defer func() {
		log.Println("[service_user.UserByID] end", "data", map[string]interface{}{"user: ": c}, "error", err)
	}()

	if err = c.Valid(); err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.InvalidErr
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf(codeErr)
	}

	result = &model_user.User{}
	err = collection.User().Collection().FindOne(ctx, bson.M{"_id": c.UserID}).Decode(result)
	if err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf(codeErr)
	}

	return result, nil
}
