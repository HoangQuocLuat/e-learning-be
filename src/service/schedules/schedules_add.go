package service_schedules

import (
	"context"
	"fmt"
	"time"

	"e-learning/src/database/collection"

	"go.mongodb.org/mongo-driver/bson/primitive"

	src_const "e-learning/src/const"
	model_schedules "e-learning/src/database/model/schedules"
)

type SchedulesAddCommand struct {
	ClassID       string
	SchedulesType string
	Description   string
	DayOfWeek     string
	StartDate     string
	EndDate       string
	StartTime     string
	EndTime       string
}

func SchedulesAdd(ctx context.Context, c *SchedulesAddCommand) (result *model_schedules.Schedules, err error) {

	startDate, err := time.Parse("02-01-2006", c.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid StartDate format: %v", err)
	}

	endDate, err := time.Parse("02-01-2006", c.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid EndDate format: %v", err)
	}

	startTime, err := time.Parse("15:04", c.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid StartTime format: %v", err)
	}

	endTime, err := time.Parse("15:04", c.EndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid EndTime format: %v", err)
	}
	result = &model_schedules.Schedules{
		ID:            primitive.NewObjectID().Hex(),
		ClassID:       c.ClassID,
		Description:   c.Description,
		SchedulesType: c.SchedulesType,
		StartDate:     startDate,
		EndDate:       endDate,
		StartTime:     startTime,
		EndTime:       endTime,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err = collection.Schedules().Collection().InsertOne(ctx, result)
	if err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Schedules + src_const.InternalError
		return nil, fmt.Errorf(codeErr)
	}
	return
}
