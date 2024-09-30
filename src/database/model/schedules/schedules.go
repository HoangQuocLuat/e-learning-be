package model_schedules

import (
	"time"

	graphModel "e-learning/src/graph/generated/model"
)


type Schedules struct {
	ID            string    `json:"id" bson:"_id"`
	ClassID       string    `json:"class_id" bson:"class_id"`
	ClassName     string    `json:"class_name" bson:"class_name"`
	Description   string    `json:"description" bson:"description"`
	SchedulesType string    `json:"schedules_type" bson:"schedules_type"`
	DayOfWeek     int    `json:"day_of_week" bson:"day_of_week"`
	StartTime     time.Time `json:"start_time" bson:"start_time"`
	EndTime       time.Time `json:"end_time" bson:"end_time"`
	StartDate     time.Time `json:"start_date" bson:"start_date"`
	EndDate       time.Time `json:"end_date" bson:"end_date"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
}

func (a *Schedules) ConvertToModelGraph() *graphModel.Schedules {
	data := graphModel.Schedules{
		ID:            a.ID,
		Description:   a.Description,
		SchedulesType: a.SchedulesType,
		DayOfWeek:     a.DayOfWeek,
		StartDate:     a.StartDate,
		EndDate:       a.EndDate,
		StartTime:     a.StartTime,
		EndTime:       a.EndTime,
		Class: &graphModel.Class{
			ID:        a.ClassID,
			ClassName: a.ClassName,
		},
	}

	return &data
}
