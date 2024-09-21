package service_auth

import (
	"context"
	"e-learning/config"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	model_user "e-learning/src/database/model/user"
	"e-learning/src/service"
	"e-learning/src/utilities"
	"fmt"

	"github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type AuthLoginCommand struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func (c *AuthLoginCommand) Valid() error {
	_, err := govalidator.ValidateStruct(c)
	return err
}

func AuthBasicLogin(ctx context.Context, c *AuthLoginCommand) (accessToken string, refreshToken string, role string, err error) {
	if err = c.Valid(); err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.InvalidErr
		service.AddError(ctx, "", "", codeErr)
		return "", "", "", err
	}

	condition := make(map[string]interface{})
	condition["user_name"] = c.UserName

	account := &model_user.User{}
	err = collection.User().Collection().FindOne(ctx, bson.M{"user_name": c.UserName}).Decode(account)
	if err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.UserNotFound
		service.AddError(ctx, "", "", codeErr)
		return "", "", "", fmt.Errorf("account is not existed")
	}
	role = account.Role

	if account.Status != model_user.StatusActive {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.UserNotActive
		service.AddError(ctx, "", "", codeErr)
		return "", "", "", fmt.Errorf("account is not active")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(c.Password))
	if err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.Unauthorized
		service.AddError(ctx, "", "", codeErr)
		return "", "", "", fmt.Errorf("unauthorized: %s", err.Error())
	}

	accessToken, err = utilities.CreateToken(account.ID, account.UserName, account.Role, account.Status, int(config.Get().ExpiresTimeAccessToken), config.Get().AccessTokenType)
	if err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.UnableCreateToken
		service.AddError(ctx, "", "", codeErr)
		return accessToken, refreshToken, "", fmt.Errorf("unable to create access token: %s", err.Error())
	}

	refreshToken, err = utilities.CreateToken(account.ID, account.UserName, account.Role, account.Status, int(config.Get().ExpiresTimeRefreshToken), config.Get().RefreshTokenType)
	if err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.UnableCreateToken
		service.AddError(ctx, "", "", codeErr)
		return accessToken, refreshToken, "", fmt.Errorf("unable to create refresh token: %s", err.Error())
	}

	return accessToken, refreshToken, role, nil
}
