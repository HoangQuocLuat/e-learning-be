package model_account

import (
	"time"

	graphModel "e-learning/src/graph/generated/model"
)

const (
	StatusActive = 1
	StatusBlock  = 2
)

const (
	RoleSuperAdmin = "super-admin"
	RoleAdmin      = "admin"
	RoleUser       = "user"
)

type Account struct {
	ID string `json:"id" bson:"_id"`

	UserName string `json:"user_name" bson:"user_name"`
	Password string `json:"password" bson:"password"`
	Status   int    `json:"status" bson:"status"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`

	LogsStatus []LogStatus `json:"logs_status" bson:"logs_status,omitempty"`

	Role string `json:"role" bson:"role"`
}

type LogStatus struct {
	Status    int32     `json:"status"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

func (a *Account) ConvertToModelGraph() *graphModel.Account {
	data := graphModel.Account{
		ID: a.ID,

		Role:   a.Role,
		Status: a.Status,
	}

	return &data
}
