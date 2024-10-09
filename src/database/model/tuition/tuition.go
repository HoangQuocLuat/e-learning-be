package model_tuition

import (
	"time"

	graphModel "e-learning/src/graph/generated/model"
)

type Tuition struct {
	ID           string    `json:"id" bson:"_id"`
	UserID       string    `json:"user_id" bson:"user_id"`
	Status       string    `json:"status" bson:"status"`
	Price        float64   `json:"price" bson:"price"`
	LessonsCount int       `json:"lessons_count" bson:"lessons_count"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}

func (a *Tuition) ConvertToModelGraph() *graphModel.Tuition {
	data := graphModel.Tuition{
		ID:           a.ID,
		UserID:       a.UserID,
		Price:        a.Price,
		Status:       a.Status,
		LessonsCount: a.LessonsCount,
	}

	return &data
}
