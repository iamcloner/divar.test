package mongodb

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	connection *Handler
	once       sync.Once
)

type Handler struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func GetMongoDBHandler() (*Handler, error) {
	var err error
	once.Do(func() {
		var client *mongo.Client
		clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_ADDRESS"))
		client, err = mongo.Connect(context.TODO(), clientOptions)
		connection = &Handler{
			Client:   client,
			Database: client.Database(os.Getenv("MONGO_DBNAME")),
		}
	})

	return connection, err
}
func (handler *Handler) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return handler.Client.Disconnect(ctx)
}

func (handler *Handler) FindOne(collectionName string, filter interface{}, projection interface{}) *mongo.SingleResult {

	collection := handler.Database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findOptions := options.FindOne().SetProjection(projection)
	result := collection.FindOne(ctx, filter, findOptions)

	return result
}
func (handler *Handler) FindMany(collectionName string, filter interface{}, projection interface{}) ([]bson.M, error) {
	collection := handler.Database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find().SetProjection(projection)
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Println(err)
		}
	}(cursor, ctx)

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

func (handler *Handler) Insert(collectionName string, document interface{}) (*mongo.InsertOneResult, error) {
	collection := handler.Database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (handler *Handler) Update(collectionName string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	collection := handler.Database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (handler *Handler) Delete(collectionName string, filter interface{}) (*mongo.DeleteResult, error) {
	collection := handler.Database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}
