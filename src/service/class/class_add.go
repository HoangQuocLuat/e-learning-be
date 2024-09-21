package service_class

import (
	"context"
	"fmt"
	"time"

	"e-learning/src/database/collection"

	"go.mongodb.org/mongo-driver/bson/primitive"

	src_const "e-learning/src/const"
	model_class "e-learning/src/database/model/class"
)

type ClassAddCommand struct {
	UserID      string
	SchedulesID string
	ClassName   string
}

func (c *ClassAddCommand) Valid() error {
	if c.ClassName == "" {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Class + src_const.InvalidErr
		return fmt.Errorf(codeErr)
	}

	return nil
}

func ClassAdd(ctx context.Context, c *ClassAddCommand) (result *model_class.Class, err error) {
	if err = c.Valid(); err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Class + src_const.InvalidErr
		return nil, fmt.Errorf(codeErr)
	}

	condition := make(map[string]interface{})
	condition["class_name"] = c.ClassName

	cnt, err := collection.Class().Collection().CountDocuments(ctx, condition)
	if err == nil && cnt > 0 {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Class + src_const.ClassExist
		return nil, fmt.Errorf(codeErr)
	}

	result = &model_class.Class{
		ID:          primitive.NewObjectID().Hex(),
		UserID:      c.UserID,
		SchedulesID: c.SchedulesID,
		ClassName:   c.ClassName,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = collection.Class().Collection().InsertOne(ctx, result)
	if err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Class + src_const.InternalError
		return nil, fmt.Errorf(codeErr)
	}
	return
}
