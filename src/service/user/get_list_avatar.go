package service_user

import (
	"context"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	model_user "e-learning/src/database/model/user"
	"e-learning/src/service"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func GetImagesDescByClassID(ctx context.Context, classID string) (res []*model_user.User, err error) {
	cursor, err := collection.User().Collection().Find(ctx, bson.M{"class_id": classID})
	if err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf("%s: %v", codeErr, err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var users model_user.User
		err := cursor.Decode(&users)
		if err != nil {
			codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InternalError
			service.AddError(ctx, "", "", codeErr)
			return nil, fmt.Errorf(codeErr)
		}

		res = append(res, &users)
	}

	if err := cursor.Err(); err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf(codeErr)
	}

	return
}
