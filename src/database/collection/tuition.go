package collection

import (
	"fmt"
	"sync"

	mongoDriver "go.mongodb.org/mongo-driver/mongo"

	"e-learning/config"

	"e-learning/mongo"
)

const TuitionCollection = "tuition"

var (
	_tuitionCollection        *TuitionMongoCollection
	loadTuitionRepositoryOnce sync.Once
)

type TuitionMongoCollection struct {
	client         *mongo.MongoDB
	collectionName string
	databaseName   string
	indexed        map[string]bool
}

func LoadTuitionCollectionMongo(mongoClient *mongo.MongoDB) (err error) {
	loadTuitionRepositoryOnce.Do(func() {
		_tuitionCollection, err = NewTuitionMongoCollection(mongoClient, config.Get().DatabaseName)
	})
	return
}

func Tuition() *TuitionMongoCollection {
	if _tuitionCollection == nil {
		panic("database: like tuition collection is not initiated")
	}
	return _tuitionCollection
}

func NewTuitionMongoCollection(client *mongo.MongoDB, databaseName string) (*TuitionMongoCollection, error) {
	if client == nil {
		return nil, fmt.Errorf("[NewTuitionMongoCollection] client nil pointer")
	}
	repo := &TuitionMongoCollection{
		client:         client,
		collectionName: TuitionCollection,
		databaseName:   databaseName,
		indexed:        make(map[string]bool),
	}
	return repo, nil
}

func (repo *TuitionMongoCollection) Collection() *mongoDriver.Collection {
	return repo.client.Client().Database(repo.databaseName).Collection(repo.collectionName)
}
