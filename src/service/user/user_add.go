package service_user

import (
	"context"
	"fmt"
	"time"

	"e-learning/src/database/collection"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	src_const "e-learning/src/const"
	model_user "e-learning/src/database/model/user"
)

type UserAddCommand struct {
	ClassID   string
	UserName  string
	Password  string
	Role      string
	Name      string
	DateBirth string
	Phone     string
	Email     string
	Address   string
}

func (c *UserAddCommand) Valid() error {
	if c.UserName == "" {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.InvalidErr
		return fmt.Errorf(codeErr)
	}

	return nil
}

func UserAdd(ctx context.Context, c *UserAddCommand) (result *model_user.User, err error) {
	if err = c.Valid(); err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.InvalidErr
		return nil, fmt.Errorf(codeErr)
	}

	condition := make(map[string]interface{})
	condition["user_name"] = c.UserName

	cnt, err := collection.User().Collection().CountDocuments(ctx, condition)
	if err == nil && cnt > 0 {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.AccountExist
		return nil, fmt.Errorf(codeErr)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.WrongPassword
		return nil, fmt.Errorf(codeErr)
	}

	result = &model_user.User{
		ID:        primitive.NewObjectID().Hex(),
		ClassID:   c.ClassID,
		UserName:  c.UserName,
		Password:  string(password),
		Role:      c.Role,
		Status:    model_user.StatusActive,
		Name:      c.Name,
		DateBirth: c.DateBirth,
		Phone:     c.Phone,
		Email:     c.Email,
		Address:   c.Address,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		LogsStatus: []model_user.LogStatus{
			{
				Status:    model_user.StatusActive,
				CreatedAt: time.Now(),
			},
		},
	}

	_, err = collection.User().Collection().InsertOne(ctx, result)
	if err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.InternalError
		return nil, fmt.Errorf(codeErr)
	}
	return
}
