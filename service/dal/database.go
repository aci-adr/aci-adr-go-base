package dal

import "go.mongodb.org/mongo-driver/bson"

type Database[T any] interface {
	GetOne(filter bson.D) (T, error)
	Get(filter bson.D) ([]T, error)
	CreateOne(document T) (T, error)
	UpdateOne(document bson.D, filter bson.D) (T, error)
	DeleteOne(filter bson.D) (int64, error)
}
