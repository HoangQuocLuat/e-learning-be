package service_account

import (
	"context"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	model_account "e-learning/src/database/model/account"
	"e-learning/src/service"
	"fmt"
	"log"
	"strings"

	mongoDriver "go.mongodb.org/mongo-driver/mongo"

	"github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson"
)

type AccountPaginationCommand struct {
	Page    int                    `json:"page"`
	Limit   int                    `json:"limit"`
	OrderBy string                 `json:"order_by"`
	Search  map[string]interface{} `json:"search"`
}

func (c *AccountPaginationCommand) Valid() error {
	if c.Page < 1 {
		c.Page = 1
	}

	if c.Limit < 1 {
		c.Limit = 10
	}

	_, err := govalidator.ValidateStruct(c)
	if err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_Account + src_const.InvalidErr
		return fmt.Errorf(codeErr)
	}
	return nil
}

func AccountPagination(ctx context.Context, c *AccountPaginationCommand) (total int, results []model_account.Account, err error) {
	log.Println("[service_account.AccountPagination] start")
	defer func() {
		log.Println("[service_account.AccountPagination] end", "data", map[string]interface{}{"command: ": c}, "error", err)
	}()

	if err = c.Valid(); err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_Account + src_const.InvalidErr
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

	cur, err := collection.Account().Collection().Aggregate(ctx, pipeline)
	if err != nil {
		codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_Account + src_const.InternalError
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
	accounts := make([]model_account.Account, 0)

	if len(listOrder["total"].(bson.A)) > 0 {
		total = int(listOrder["total"].(bson.A)[0].(bson.M)["count"].(int32))
		rows := listOrder["rows"].(bson.A)

		for _, rawAccount := range rows {
			accountBSON, err := bson.Marshal(rawAccount)
			if err != nil {
				codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_Account + src_const.InternalError
				service.AddError(ctx, "", "", codeErr)
				return 0, nil, fmt.Errorf(codeErr)
			}

			var account model_account.Account
			err = bson.Unmarshal(accountBSON, &account)
			if err != nil {
				codeErr := src_const.ServiceErr_Auth + src_const.ElementErr_Account + src_const.InternalError
				service.AddError(ctx, "", "", codeErr)
				return 0, nil, fmt.Errorf(codeErr)
			}

			accounts = append(accounts, account)
		}
	}

	return int(total), accounts, nil
}
