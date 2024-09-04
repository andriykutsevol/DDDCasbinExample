package storage

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// We use this to not depend on mongo.Collection in internal/app/app.go for example
type MongoStorage struct {
	collection *DatabaseCollection
}

func (s *MongoStorage) SetCollection(c *DatabaseCollection) {
	s.collection = c
}

func (s *MongoStorage) GetDatabaseCollection() *DatabaseCollection {
	return s.collection
}

func (s *MongoStorage) GetCollection() *mongo.Collection {
	return s.collection.collection
}

func (s *MongoStorage) AddFilterCondition(filter *bson.D, key string, value interface{}) {
	*filter = append(*filter, bson.E{Key: key, Value: value})
}
