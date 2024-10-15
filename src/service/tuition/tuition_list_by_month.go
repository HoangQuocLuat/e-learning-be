package service_tuition

import (
	"context"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	model_tuition "e-learning/src/database/model/tuition"
	model_user "e-learning/src/database/model/user"
	"e-learning/src/service"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TuitionGetListCommand struct {
	Year  string
	Month string
}

func TuitionGetList(ctx context.Context, s *TuitionGetListCommand) (results []*model_tuition.Tuition, err error) {
	startDate, _ := time.Parse("2006-01-02", fmt.Sprintf("%s-%s-01", s.Year, s.Month)) // ngày đầu tiên của tháng
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)
	req := bson.M{
		"created_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}
	cursor, err := collection.Tuition().Collection().Find(ctx, req)
	if err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Tuition + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return results, fmt.Errorf(codeErr)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var tuition model_tuition.Tuition
		err := cursor.Decode(&tuition)
		if err != nil {
			codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Tuition + src_const.InternalError
			service.AddError(ctx, "", "", codeErr)
			return results, fmt.Errorf(codeErr)
		}

		user, err := GetUserByID(ctx, tuition.UserID)
		if err != nil {
			codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InternalError
			service.AddError(ctx, "", "", codeErr)
			return results, fmt.Errorf(codeErr)
		}
		if user != nil {
			tuition.Name = user.Name
		}
		results = append(results, &tuition)
	}

	if err := cursor.Err(); err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Tuition + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return results, fmt.Errorf(codeErr)
	}

	return results, nil
}

func GetUserByID(ctx context.Context, userID string) (user *model_user.User, err error) {
	err = collection.User().Collection().FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error retrieving class: %v", err)
	}
	return user, nil
}
