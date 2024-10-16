package service_tuition

import (
	"context"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	model_tuition "e-learning/src/database/model/tuition"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type TuitionUpdateCommand struct {
	ID           string
	TotalFee     *int
	Discount     *int
	PaidAmount   *int
	RemainingFee *int
}

func (t *TuitionUpdateCommand) Valid() error {
	if t.ID == "" {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Tuition + src_const.InvalidErr
		return fmt.Errorf(codeErr)
	}
	return nil
}
func TuitionUpdate(ctx context.Context, t *TuitionUpdateCommand) (result *model_tuition.Tuition, err error) {
	if err := t.Valid(); err != nil {
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Tuition + src_const.InvalidErr
		return nil, fmt.Errorf(codeErr)
	}

	err = collection.Tuition().Collection().FindOne(ctx, bson.M{"_id": t.ID}).Decode(&result)

	if err != nil {
		log.Println("[service_tuition.TuitionUpdate]", "FindOne ID", map[string]interface{}{"command: ": t}, "error", err)
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Tuition + src_const.TuitionExist
		return nil, fmt.Errorf(codeErr)
	}

	updateTuition := bson.M{}

	if t.Discount != nil {
		updateTuition["discount"] = *t.Discount
		result.Discount = *t.Discount
	}
	if t.TotalFee != nil {
		updateTuition["total_fee"] = *t.TotalFee
		result.TotalFee = *t.TotalFee
	}
	if t.PaidAmount != nil {
		updateTuition["paid_amount"] = *t.PaidAmount
		result.PaidAmount = *t.PaidAmount
	}
	if t.RemainingFee != nil {
		updateTuition["remaining_fee"] = *t.RemainingFee
		result.RemainingFee = *t.RemainingFee
	}

	if t.Discount != nil || t.TotalFee != nil || t.PaidAmount != nil || t.RemainingFee != nil {
		updateTuition["updated_at"] = time.Now()
		result.UpdatedAt = time.Now()
	}

	_, err = collection.Tuition().Collection().UpdateOne(ctx, bson.M{"_id": t.ID}, bson.M{"$set": updateTuition})

	if err != nil {
		log.Println("[service_order.TuitionUpdate]", "Update", map[string]interface{}{"command: ": t}, "error", err)
		codeErr := src_const.ServiceErr_E_Learning + src_const.ElementErr_Tuition + src_const.InternalError
		return nil, fmt.Errorf(codeErr)
	}

	return
}
