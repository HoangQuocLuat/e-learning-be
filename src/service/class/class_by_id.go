package service_class

import (
	"context"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	model_class "e-learning/src/database/model/class"
	"e-learning/src/service"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type ClassByIDCommand struct {
	ClassID string
}

func (c *ClassByIDCommand) Valid() error {
	if c.ClassID == "" {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Class + src_const.InvalidErr
		return fmt.Errorf(codeErr)
	}

	return nil
}

func ClassByID(ctx context.Context, c *ClassByIDCommand) (result *model_class.Class, err error) {
	log.Println("[service_class.ClassByID] start")
	defer func() {
		log.Println("[service_class.ClassByID] end", "data", map[string]interface{}{"class: ": c}, "error", err)
	}()

	if err = c.Valid(); err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Class + src_const.InvalidErr
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf(codeErr)
	}

	result = &model_class.Class{}
	err = collection.Class().Collection().FindOne(ctx, bson.M{"_id": c.ClassID}).Decode(result)
	if err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Class + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf(codeErr)
	}

	return result, nil
}
