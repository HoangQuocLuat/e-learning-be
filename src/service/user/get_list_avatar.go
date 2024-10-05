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

func GetImagesByClassID(ctx context.Context, classID string) ([]string, error) {
	cursor, err := collection.User().Collection().Find(ctx, bson.M{"class_id": classID})
	if err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf("%s: %v", codeErr, err)
	}
	defer cursor.Close(ctx)

	var images []string
	for cursor.Next(ctx) {
		var user model_user.User
		if err := cursor.Decode(&user); err != nil {
			codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InternalError
			service.AddError(ctx, "", "", codeErr)
			return nil, fmt.Errorf("%s: %v", codeErr, err)
		}
		images = append(images, user.Avatar)
	}

	if err := cursor.Err(); err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf("%s: %v", codeErr, err)
	}
	return images, nil
}
