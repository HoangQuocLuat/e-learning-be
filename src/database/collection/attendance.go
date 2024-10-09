package collection

import (
	"fmt"
	"sync"

	mongoDriver "go.mongodb.org/mongo-driver/mongo"

	"e-learning/config"

	"e-learning/mongo"
)

const AttendanceCollection = "attendance"

var (
	_attendanceCollection        *AttendanceMongoCollection
	loadAttendanceRepositoryOnce sync.Once
)

type AttendanceMongoCollection struct {
	client         *mongo.MongoDB
	collectionName string
	databaseName   string
	indexed        map[string]bool
}

func LoadAttendanceCollectionMongo(mongoClient *mongo.MongoDB) (err error) {
	loadAttendanceRepositoryOnce.Do(func() {
		_attendanceCollection, err = NewAttendanceMongoCollection(mongoClient, config.Get().DatabaseName)
	})
	return
}

func Attendance() *AttendanceMongoCollection {
	if _attendanceCollection == nil {
		panic("database: like attendance collection is not initiated")
	}
	return _attendanceCollection
}

func NewAttendanceMongoCollection(client *mongo.MongoDB, databaseName string) (*AttendanceMongoCollection, error) {
	if client == nil {
		return nil, fmt.Errorf("[NewAttendanceMongoCollection] client nil pointer")
	}
	repo := &AttendanceMongoCollection{
		client:         client,
		collectionName: AttendanceCollection,
		databaseName:   databaseName,
		indexed:        make(map[string]bool),
	}
	return repo, nil
}

func (repo *AttendanceMongoCollection) Collection() *mongoDriver.Collection {
	return repo.client.Client().Database(repo.databaseName).Collection(repo.collectionName)
}
