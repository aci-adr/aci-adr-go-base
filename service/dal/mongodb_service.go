package dal

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDbService[T any] struct {
	Collection string
}

var database *mongo.Database

func InitMongo() {
	loggerOptions := options.
		Logger().
		SetComponentLevel(options.LogComponentServerSelection, options.LogLevelDebug)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetLoggerOptions(loggerOptions))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	database = client.Database(os.Getenv("DB_NAME"))
}

func getDatabase() *mongo.Database {
	return database
}

func getContext() context.Context {
	return context.Background()
}

func (db *MongoDbService[T]) GetOne(filter bson.D) (T, error) {
	result := getDatabase().Collection(db.Collection).FindOne(getContext(), filter)
	var data T
	err := result.Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (db *MongoDbService[T]) Get(filter bson.D) ([]T, error) {
	cursor, _ := getDatabase().Collection(db.Collection).Find(getContext(), filter)
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, context.Background())
	var data []T
	for cursor.Next(context.Background()) {
		var result T
		if err := cursor.Decode(&result); err != nil {
			fmt.Println("Error decoding document:", err)
		}
		data = append(data, result)
	}
	return data, nil
}

func (db *MongoDbService[T]) CreateOne(document T) (T, error) {
	_, err := getDatabase().Collection(db.Collection).InsertOne(getContext(), document)

	if err != nil {
		return document, err
	}
	return document, nil
}

func (db *MongoDbService[T]) UpdateOne(document bson.D, filter bson.D) (T, error) {
	option := options.FindOneAndUpdate().SetReturnDocument(options.After)

	result := getDatabase().Collection(db.Collection).FindOneAndUpdate(getContext(), filter, document, option)

	var data T
	err := result.Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (db *MongoDbService[T]) DeleteOne(filter bson.D) (int64, error) {
	result, err := getDatabase().Collection(db.Collection).DeleteOne(getContext(), filter)

	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
