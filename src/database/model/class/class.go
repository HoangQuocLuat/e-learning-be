package model_class

import (
	"time"

	graphModel "e-learning/src/graph/generated/model"
)

type Class struct {
	ID          string    `json:"id" bson:"_id"`
	UserID      string    `json:"user_id" bson:"user_id"`
	SchedulesID string    `json:"schedules_id" bson:"schedules_id"`
	ClassName   string    `json:"class_name" bson:"class_name"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

func (a *Class) ConvertToModelGraph() *graphModel.Class {
	data := graphModel.Class{
		ID:        a.ID,
		ClassName: a.ClassName,
	}

	return &data
}
