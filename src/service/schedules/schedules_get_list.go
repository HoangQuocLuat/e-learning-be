package service_schedules

import (
	"context"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	model_class "e-learning/src/database/model/class" // Thêm import cho model class
	model_schedules "e-learning/src/database/model/schedules"
	"e-learning/src/service"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SchedulesGetListCommand struct{}

// Hàm để lấy thông tin class theo class_id
func GetClassByID(ctx context.Context, classID string) (*model_class.Class, error) {
	var class model_class.Class
	err := collection.Class().Collection().FindOne(ctx, bson.M{"_id": classID}).Decode(&class)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Trường hợp không tìm thấy class với class_id này
			return nil, nil
		}
		return nil, fmt.Errorf("error retrieving class: %v", err)
	}
	return &class, nil
}

// Hàm lấy danh sách Schedules kèm thông tin Class
func SchedulesGetList(ctx context.Context, s *SchedulesGetListCommand) (results []*model_schedules.Schedules, err error) {
	cursor, err := collection.Schedules().Collection().Find(ctx, bson.M{})
	if err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return results, fmt.Errorf(codeErr)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var schedules model_schedules.Schedules
		err := cursor.Decode(&schedules)
		if err != nil {
			codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InternalError
			service.AddError(ctx, "", "", codeErr)
			return results, fmt.Errorf(codeErr)
		}

		class, err := GetClassByID(ctx, schedules.ClassID)
		if err != nil {
			codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InternalError
			service.AddError(ctx, "", "", codeErr)
			return results, fmt.Errorf(codeErr)
		}

		if class != nil {
			schedules.ClassName = class.ClassName
		}

		results = append(results, &schedules)
	}

	if err := cursor.Err(); err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return results, fmt.Errorf(codeErr)
	}

	return results, nil
}
