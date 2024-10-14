package model_order

import (
	graphModel "e-learning/src/graph/generated/model"
)

type Order struct {
	OrderURL string `json:"order_url"`
}

func (a *Order) ConvertToModelGraph() *graphModel.Order {
	data := graphModel.Order{
		OrderURL: a.OrderURL,
	}

	return &data
}
