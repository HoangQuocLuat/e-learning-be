package service_account

import (
	"context"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	"e-learning/src/service"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type AccountDeleteCommand struct {
	ID string
}

func (c *AccountDeleteCommand) Valid() error {
	if c.ID == "" {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_Account + src_const.InvalidErr
		return fmt.Errorf(codeErr)
	}

	return nil
}

func AccountDelete(ctx context.Context, c *AccountDeleteCommand) (err error) {
	log.Println("[service_account.AccountDelete] start")
	defer func() {
		log.Println("[service_account.AccountDelete] end", "data", map[string]interface{}{"account: ": c}, "error", err)
	}()

	if err = c.Valid(); err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_Account + src_const.InvalidErr
		service.AddError(ctx, "", "", codeErr)
		return fmt.Errorf(codeErr)
	}
	_, err = collection.Account().Collection().DeleteOne(ctx, bson.M{"_id": c.ID})
	if err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_Account + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return fmt.Errorf(codeErr)
	}
	return
}
