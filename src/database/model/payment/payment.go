package model_payment

import (
	graphModel "e-learning/src/graph/generated/model"
	"time"
)

type Payment struct {
	ID        string    `json:"id" bson:"_id"`
	UserID    string    `json:"user_id" bson:"user_id"`
	Name      string    `json:"name"`
	TuitionID string    `json:"tuition_id" bson:"tuition_id"`
	Amount    string    `json:"amount" bson:"amount"`
	TransID   string    `json:"trans_id" bson:"trans_id"`
	Status    string    `json:"status" bson:"status"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func (a *Payment) ConvertToModelGraph() *graphModel.Payment {
	location, _ := time.LoadLocation("Asia/Bangkok")
	m := a.UpdatedAt.In(location).Format("2006-01-02 15:04:05")
	data := graphModel.Payment{
		ID:      a.ID,
		Amount:  a.Amount,
		TransID: a.TransID,
		Status:  a.Status,
		Date:    &m,
		User: graphModel.User{
			ID:   a.UserID,
			Name: a.Name,
		},
		Tuition: graphModel.Tuition{
			ID: a.TuitionID,
		},
	}

	return &data
}
