package service_attendence

import (
	"context"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	model_attendance "e-learning/src/database/model/attendance"
	"e-learning/src/service"
	service_payment "e-learning/src/service/payment"
	"fmt"
	"log"
	"strings"
	"time"

	mongoDriver "go.mongodb.org/mongo-driver/mongo"

	"github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson"
)

type AttendancePaginationCommand struct {
	Page    int                    `json:"page"`
	Limit   int                    `json:"limit"`
	OrderBy string                 `json:"order_by"`
	Search  map[string]interface{} `json:"search"`
}

func (c *AttendancePaginationCommand) Valid() error {
	if c.Page < 1 {
		c.Page = 1
	}

	if c.Limit < 1 {
		c.Limit = 10
	}

	_, err := govalidator.ValidateStruct(c)
	if err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.InvalidErr
		return fmt.Errorf(codeErr)
	}
	return nil
}

func AttendancePagination(ctx context.Context, c *AttendancePaginationCommand) (total int, results []model_attendance.Attendance, err error) {
	log.Println("[service_user.UserPagination] start")
	defer func() {
		log.Println("[service_user.UserPagination] end", "data", map[string]interface{}{"command: ": c}, "error", err)
	}()

	if err = c.Valid(); err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.InvalidErr
		service.AddError(ctx, "", "", codeErr)
		return 0, nil, fmt.Errorf(codeErr)
	}
	condition := make(map[string]interface{})

	if phone, ok := c.Search["phone"]; ok {
		condition["phone"] = phone
	}

	if user_name, ok := c.Search["user_name"]; ok {
		condition["user_name"] = user_name
	}

	objOrderBy := bson.M{}
	if c.OrderBy != "" {
		value := src_const.ASC
		if strings.HasPrefix(c.OrderBy, "-") {
			value = src_const.DESC
			c.OrderBy = strings.TrimPrefix(c.OrderBy, "-")
		}

		objOrderBy = bson.M{c.OrderBy: value}
	}

	//Default order by updated_at | new -> old
	if c.OrderBy == "" {
		objOrderBy = bson.M{"updated_at": src_const.DESC}
	}

	matchStage := bson.D{{Key: "$match", Value: condition}}

	facectStage := bson.D{{
		Key: "$facet",

		Value: bson.M{
			"rows": bson.A{
				bson.M{"$skip": (c.Page - 1) * c.Limit},
				bson.M{"$limit": c.Limit},
			},
			"total": bson.A{
				bson.M{"$count": "count"},
			},
		},
	}}

	sortStage := bson.D{{Key: "$sort", Value: objOrderBy}}

	pipeline := mongoDriver.Pipeline{
		matchStage,
		sortStage,
		facectStage,
	}

	cur, err := collection.Attendance().Collection().Aggregate(ctx, pipeline)
	if err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.InternalError
		return 0, nil, fmt.Errorf(codeErr)
	}

	var listOrder bson.M
	for cur.Next(ctx) {
		err := cur.Decode(&listOrder)
		if err != nil {
			codeErr := src_const.ServiceErr_Auth + src_const.ServiceErr_Auth + src_const.InternalError
			service.AddError(ctx, "", "", codeErr)
			return 0, nil, fmt.Errorf(codeErr)
		}
	}

	// Extract the total count and rows from the result
	ats := make([]model_attendance.Attendance, 0)

	if len(listOrder["total"].(bson.A)) > 0 {
		total = int(listOrder["total"].(bson.A)[0].(bson.M)["count"].(int32))
		rows := listOrder["rows"].(bson.A)

		for _, rawUser := range rows {
			userBSON, err := bson.Marshal(rawUser)
			if err != nil {
				codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.InternalError
				service.AddError(ctx, "", "", codeErr)
				return 0, nil, fmt.Errorf(codeErr)
			}

			var attendance model_attendance.Attendance
			err = bson.Unmarshal(userBSON, &attendance)
			if err != nil {
				codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_User + src_const.InternalError
				service.AddError(ctx, "", "", codeErr)
				return 0, nil, fmt.Errorf(codeErr)
			}
			location, err := time.LoadLocation("Asia/Bangkok")
			if err != nil {
				log.Fatalf("Lỗi khi tải múi giờ: %v", err)
			}

			attendance.TimeCheckIn = attendance.TimeCheckIn.In(location)

			a, err := service_payment.GetNameByID(ctx, attendance.ID)
			if err != nil {
				codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InternalError
				service.AddError(ctx, "", "", codeErr)
				return 0, nil, fmt.Errorf(codeErr)
			}
			if a != nil {
				attendance.Name = a.Name
			}

			ats = append(ats, attendance)
		}
	}

	return int(total), ats, nil
}
