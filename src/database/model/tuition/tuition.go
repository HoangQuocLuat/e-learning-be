package model_tuition

import (
	"time"

	graphModel "e-learning/src/graph/generated/model"
)

type Tuition struct {
	ID           string    `json:"id" bson:"_id"`
	UserID       string    `json:"user_id" bson:"user_id"`
	Name         string    `json:"name" bson:"name"`
	TotalFee     int       `json:"total_fee" bson:"total_fee"`
	Discount     int       `json:"discount" bson:"discount"`
	PaidAmount   int       `json:"paid_amount" bson:"paid_amount"`
	RemainingFee int       `json:"remaining_fee" bson:"remaining_fee"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}

func (a *Tuition) ConvertToModelGraph() *graphModel.Tuition {
	data := graphModel.Tuition{
		ID:           a.ID,
		TotalFee:     a.TotalFee,
		Discount:     &a.Discount,
		PaidAmount:   &a.PaidAmount,
		RemainingFee: &a.RemainingFee,
		User: &graphModel.User{
			ID:   a.UserID,
			Name: a.Name,
		},
	}

	return &data
}
