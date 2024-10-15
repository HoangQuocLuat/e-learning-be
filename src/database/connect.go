package database

import (
	"context"
	"e-learning/config"
	"e-learning/mongo"
	src_const "e-learning/src/const"
	"e-learning/src/database/collection"
	"log"
	"time"
)

func ConnectDatabse(ctx context.Context) error {
	var mongoClient *mongo.MongoDB
	var err error
	numberRetry := config.Get().NumberRetry
	if numberRetry == 0 {
		numberRetry = src_const.DEFAULTNUMBERRETRY
	}

	for i := 1; i <= config.Get().NumberRetry; i++ {
		mongoClient, err = mongo.NewMongoDBFromUrl(ctx, config.Get().MongoURL, time.Second*10)
		if err != nil {
			if i == config.Get().NumberRetry {
				log.Println(err)
				return err
			}
			time.Sleep(10 * time.Second)
		}

		if mongoClient != nil {
			break
		}
	}

	if err := collection.LoadUserCollectionMongo(mongoClient); err != nil {
		return err
	}

	if err := collection.LoadClassCollectionMongo(mongoClient); err != nil {
		return err
	}

	if err := collection.LoadSchedulesCollectionMongo(mongoClient); err != nil {
		return err
	}

	if err := collection.LoadAttendanceCollectionMongo(mongoClient); err != nil {
		return err
	}

	if err := collection.LoadTuitionCollectionMongo(mongoClient); err != nil {
		return err
	}
	
	if err := collection.LoadPaymentCollectionMongo(mongoClient); err!= nil {
        return err
    }

	return nil
}
