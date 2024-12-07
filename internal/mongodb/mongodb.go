package mongodb

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var connection *MongoDBHandler

type MongoDBHandler struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func Init_Mongo(uri, dbName string) error {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	connection = &MongoDBHandler{
		Client:   client,
		Database: client.Database(dbName),
	}
	return nil
}

func NewMongoDBHandler() (*MongoDBHandler, error) {
	if connection == nil {
		err := Init_Mongo(os.Getenv("MONGO_ADDRESS"), os.Getenv("MONGO_DBNAME"))
		if err != nil {
			fmt.Println("Error : Can't init mongodb connection", err.Error())
		}
	}

	return connection, nil
}
func (handler *MongoDBHandler) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return handler.Client.Disconnect(ctx)
}

func (handler *MongoDBHandler) FindOne(collectionName string, filter interface{}, projection interface{}) *mongo.SingleResult {
	collection := handler.Database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findOptions := options.FindOne().SetProjection(projection)
	result := collection.FindOne(ctx, filter, findOptions)

	return result
}
func (handler *MongoDBHandler) FindMany(collectionName string, filter interface{}, projection interface{}) ([]bson.M, error) {
	collection := handler.Database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find().SetProjection(projection)
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	for cursor.Next(ctx) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (handler *MongoDBHandler) Insert(collectionName string, document interface{}) (*mongo.InsertOneResult, error) {
	collection := handler.Database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (handler *MongoDBHandler) Update(collectionName string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	collection := handler.Database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (handler *MongoDBHandler) Delete(collectionName string, filter interface{}) (*mongo.DeleteResult, error) {
	collection := handler.Database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}
