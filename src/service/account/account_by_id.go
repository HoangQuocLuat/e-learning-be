package service_account

import (
	"context"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	model_account "e-learning/src/database/model/account"
	"e-learning/src/service"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type AccountByIDCommand struct {
	AccountID string
}

func (c *AccountByIDCommand) Valid() error {
	if c.AccountID == "" {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_Account + src_const.InvalidErr
		return fmt.Errorf(codeErr)
	}

	return nil
}

func AccountByID(ctx context.Context, c *AccountByIDCommand) (result *model_account.Account, err error) {
	log.Println("[service_account.AccountByID] start")
	defer func() {
		log.Println("[service_account.AccountByID] end", "data", map[string]interface{}{"account: ": c}, "error", err)
	}()

	if err = c.Valid(); err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_Account + src_const.InvalidErr
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf(codeErr)
	}

	result = &model_account.Account{}
	err = collection.Account().Collection().FindOne(ctx, bson.M{"_id": c.AccountID}).Decode(result)
	if err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_Account + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf(codeErr)
	}

	return result, nil
}
