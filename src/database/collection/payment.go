package collection

import (
	"fmt"
	"sync"

	mongoDriver "go.mongodb.org/mongo-driver/mongo"

	"e-learning/config"

	"e-learning/mongo"
)

const PaymentCollection = "payment"

var (
	_paymentCollection        *PaymentMongoCollection
	loadPaymentRepositoryOnce sync.Once
)

type PaymentMongoCollection struct {
	client         *mongo.MongoDB
	collectionName string
	databaseName   string
	indexed        map[string]bool
}

func LoadPaymentCollectionMongo(mongoClient *mongo.MongoDB) (err error) {
	loadPaymentRepositoryOnce.Do(func() {
		_paymentCollection, err = NewPaymentMongoCollection(mongoClient, config.Get().DatabaseName)
	})
	return
}

func Payment() *PaymentMongoCollection {
	if _paymentCollection == nil {
		panic("database: like payment collection is not initiated")
	}
	return _paymentCollection
}

func NewPaymentMongoCollection(client *mongo.MongoDB, databaseName string) (*PaymentMongoCollection, error) {
	if client == nil {
		return nil, fmt.Errorf("[NewPaymentMongoCollection] client nil pointer")
	}
	repo := &PaymentMongoCollection{
		client:         client,
		collectionName: PaymentCollection,
		databaseName:   databaseName,
		indexed:        make(map[string]bool),
	}
	return repo, nil
}

func (repo *PaymentMongoCollection) Collection() *mongoDriver.Collection {
	return repo.client.Client().Database(repo.databaseName).Collection(repo.collectionName)
}
