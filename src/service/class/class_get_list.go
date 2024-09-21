package service_class

import (
	"context"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	model_class "e-learning/src/database/model/class"
	"e-learning/src/service"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type ClassGetListCommand struct{}

func ClassGetList(ctx context.Context, s *ClassGetListCommand) (results []*model_class.Class, err error) {
	cursor, err := collection.Class().Collection().Find(ctx, bson.M{})
	if err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Class + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return results, fmt.Errorf(codeErr)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var class model_class.Class
		err := cursor.Decode(&class)
		if err != nil {
			codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Class + src_const.InternalError
			service.AddError(ctx, "", "", codeErr)
			return results, fmt.Errorf(codeErr)
		}
		results = append(results, &class)
	}

	if err := cursor.Err(); err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Class + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return results, fmt.Errorf(codeErr)
	}

	return results, nil
}
