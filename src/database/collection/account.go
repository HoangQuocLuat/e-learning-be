package collection

import (
	"fmt"
	"sync"

	mongoDriver "go.mongodb.org/mongo-driver/mongo"

	"e-learning/config"

	"e-learning/mongo"
)

const AccountCollection = "account"

var (
	_accountCollection        *AccountMongoCollection
	loadAccountRepositoryOnce sync.Once
)

type AccountMongoCollection struct {
	client         *mongo.MongoDB
	collectionName string
	databaseName   string
	indexed        map[string]bool
}

func LoadAccountCollectionMongo(mongoClient *mongo.MongoDB) (err error) {
	loadAccountRepositoryOnce.Do(func() {
		_accountCollection, err = NewAccountMongoCollection(mongoClient, config.Get().DatabaseName)
	})
	return
}

func Account() *AccountMongoCollection {
	if _accountCollection == nil {
		panic("database: like account collection is not initiated")
	}
	return _accountCollection
}

func NewAccountMongoCollection(client *mongo.MongoDB, databaseName string) (*AccountMongoCollection, error) {
	if client == nil {
		return nil, fmt.Errorf("[NewAccountMongoCollection] client nil pointer")
	}
	repo := &AccountMongoCollection{
		client:         client,
		collectionName: AccountCollection,
		databaseName:   databaseName,
		indexed:        make(map[string]bool),
	}
	return repo, nil
}

func (repo *AccountMongoCollection) Collection() *mongoDriver.Collection {
	return repo.client.Client().Database(repo.databaseName).Collection(repo.collectionName)
}
