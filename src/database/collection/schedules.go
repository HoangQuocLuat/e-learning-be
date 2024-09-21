package collection

import (
	"fmt"
	"sync"

	mongoDriver "go.mongodb.org/mongo-driver/mongo"

	"e-learning/config"

	"e-learning/mongo"
)

const SchedulesCollection = "schedules"

var (
	_schedulesCollection        *SchedulesMongoCollection
	loadSchedulesRepositoryOnce sync.Once
)

type SchedulesMongoCollection struct {
	client         *mongo.MongoDB
	collectionName string
	databaseName   string
	indexed        map[string]bool
}

func LoadSchedulesCollectionMongo(mongoClient *mongo.MongoDB) (err error) {
	loadSchedulesRepositoryOnce.Do(func() {
		_schedulesCollection, err = NewSchedulesMongoCollection(mongoClient, config.Get().DatabaseName)
	})
	return
}

func Schedules() *SchedulesMongoCollection {
	if _schedulesCollection == nil {
		panic("database: schedules collection is not initiated")
	}
	return _schedulesCollection
}

func NewSchedulesMongoCollection(client *mongo.MongoDB, databaseName string) (*SchedulesMongoCollection, error) {
	if client == nil {
		return nil, fmt.Errorf("[NewSchedulesMongoCollection] client nil pointer")
	}
	repo := &SchedulesMongoCollection{
		client:         client,
		collectionName: SchedulesCollection,
		databaseName:   databaseName,
		indexed:        make(map[string]bool),
	}
	return repo, nil
}

func (repo *SchedulesMongoCollection) Collection() *mongoDriver.Collection {
	return repo.client.Client().Database(repo.databaseName).Collection(repo.collectionName)
}
