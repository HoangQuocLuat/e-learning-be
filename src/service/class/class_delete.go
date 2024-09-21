package service_class

import (
	"context"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	"e-learning/src/service"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type ClassDeleteCommand struct {
	ID string
}

func (c *ClassDeleteCommand) Valid() error {
	if c.ID == "" {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Class + src_const.InvalidErr
		return fmt.Errorf(codeErr)
	}

	return nil
}

func ClassDelete(ctx context.Context, c *ClassDeleteCommand) (err error) {
	log.Println("[service_class.ClassDelete] start")
	defer func() {
		log.Println("[service_class.ClassDelete] end", "data", map[string]interface{}{"class: ": c}, "error", err)
	}()

	if err = c.Valid(); err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Class + src_const.InvalidErr
		service.AddError(ctx, "", "", codeErr)
		return fmt.Errorf(codeErr)
	}
	_, err = collection.Class().Collection().DeleteOne(ctx, bson.M{"_id": c.ID})
	if err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Class + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return fmt.Errorf(codeErr)
	}
	return
}
