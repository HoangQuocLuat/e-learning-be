package collection

import (
	"fmt"
	"sync"

	mongoDriver "go.mongodb.org/mongo-driver/mongo"

	"e-learning/config"

	"e-learning/mongo"
)

const ClassCollection = "class"

var (
	_classCollection        *ClassMongoCollection
	loadClassRepositoryOnce sync.Once
)

type ClassMongoCollection struct {
	client         *mongo.MongoDB
	collectionName string
	databaseName   string
	indexed        map[string]bool
}

func LoadClassCollectionMongo(mongoClient *mongo.MongoDB) (err error) {
	loadClassRepositoryOnce.Do(func() {
		_classCollection, err = NewClassMongoCollection(mongoClient, config.Get().DatabaseName)
	})
	return
}

func Class() *ClassMongoCollection {
	if _classCollection == nil {
		panic("database: like class collection is not initiated")
	}
	return _classCollection
}

func NewClassMongoCollection(client *mongo.MongoDB, databaseName string) (*ClassMongoCollection, error) {
	if client == nil {
		return nil, fmt.Errorf("[NewClassMongoCollection] client nil pointer")
	}
	repo := &ClassMongoCollection{
		client:         client,
		collectionName: ClassCollection,
		databaseName:   databaseName,
		indexed:        make(map[string]bool),
	}
	return repo, nil
}

func (repo *ClassMongoCollection) Collection() *mongoDriver.Collection {
	return repo.client.Client().Database(repo.databaseName).Collection(repo.collectionName)
}
