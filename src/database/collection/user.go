package collection

import (
	"fmt"
	"sync"

	mongoDriver "go.mongodb.org/mongo-driver/mongo"

	"e-learning/config"

	"e-learning/mongo"
)

const UserCollection = "user"

var (
	_userCollection        *UserMongoCollection
	loadUserRepositoryOnce sync.Once
)

type UserMongoCollection struct {
	client         *mongo.MongoDB
	collectionName string
	databaseName   string
	indexed        map[string]bool
}

func LoadUserCollectionMongo(mongoClient *mongo.MongoDB) (err error) {
	loadUserRepositoryOnce.Do(func() {
		_userCollection, err = NewUserMongoCollection(mongoClient, config.Get().DatabaseName)
	})
	return
}

func User() *UserMongoCollection {
	if _userCollection == nil {
		panic("database: like user collection is not initiated")
	}
	return _userCollection
}

func NewUserMongoCollection(client *mongo.MongoDB, databaseName string) (*UserMongoCollection, error) {
	if client == nil {
		return nil, fmt.Errorf("[NewUserMongoCollection] client nil pointer")
	}
	repo := &UserMongoCollection{
		client:         client,
		collectionName: UserCollection,
		databaseName:   databaseName,
		indexed:        make(map[string]bool),
	}
	return repo, nil
}

func (repo *UserMongoCollection) Collection() *mongoDriver.Collection {
	return repo.client.Client().Database(repo.databaseName).Collection(repo.collectionName)
}
