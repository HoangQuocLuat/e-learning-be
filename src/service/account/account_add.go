package service_account

import (
	"context"
	"fmt"
	"time"

	"e-learning/src/database/collection"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	src_const "e-learning/src/const"
	model_account "e-learning/src/database/model/account"
)

type AccountAddCommand struct {
	UserName string
	Password string
	Role     string
}

func (c *AccountAddCommand) Valid() error {
	if c.UserName == "" {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_Account + src_const.InvalidErr
		return fmt.Errorf(codeErr)
	}

	return nil
}

func AccountAdd(ctx context.Context, c *AccountAddCommand) (result *model_account.Account, err error) {
	if err = c.Valid(); err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_Account + src_const.InvalidErr
		return nil, fmt.Errorf(codeErr)
	}

	condition := make(map[string]interface{})
	condition["user_name"] = c.UserName

	cnt, err := collection.Account().Collection().CountDocuments(ctx, condition)
	if err == nil && cnt > 0 {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_Account + src_const.AccountExist
		return nil, fmt.Errorf(codeErr)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_Account + src_const.WrongPassword
		return nil, fmt.Errorf(codeErr)
	}

	result = &model_account.Account{
		ID: primitive.NewObjectID().Hex(),

		UserName: c.UserName,
		Password: string(password),
		Role:     c.Role,
		Status:   model_account.StatusActive,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		LogsStatus: []model_account.LogStatus{
			{
				Status:    model_account.StatusActive,
				CreatedAt: time.Now(),
			},
		},
	}

	_, err = collection.Account().Collection().InsertOne(ctx, result)
	if err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_Account + src_const.InternalError
		return nil, fmt.Errorf(codeErr)
	}
	return
}
