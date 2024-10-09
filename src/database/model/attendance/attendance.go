package model_attendance

import (
	"time"

	graphModel "e-learning/src/graph/generated/model"
)

type Attendance struct {
	ID           string    `json:"id" bson:"_id"`
	UserID       string    `json:"user_id" bson:"user_id"`
	TimeCheckIn  time.Time `json:"time_check_in" bson:"time_check_in"`
	TimeCheckOut time.Time `json:"time_check_out" bson:"time_check_out"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}

func (a *Attendance) ConvertToModelGraph() *graphModel.Attendance {
	data := graphModel.Attendance{
		ID:           a.ID,
		TimeCheckIn:  &a.TimeCheckIn,
		TimeCheckOut: &a.TimeCheckOut,
		User: graphModel.User{
			ID: a.UserID,
		},
	}

	return &data
}
