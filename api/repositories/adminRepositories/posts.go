package adminRepositories

import (
	"divar.ir/internal/mongodb"
	"divar.ir/schema"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetPendingPosts() ([]schema.Posts, error) {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return nil, errors.New("server error")
	}
	result, err := handler.FindMany("posts", bson.M{}, bson.M{})
	if err != nil {
		return nil, errors.New("server error")
	}
	resBytes, err := bson.Marshal(result)
	if err != nil {
		return nil, errors.New("server error")
	}
	var res []schema.Posts
	err = bson.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, errors.New("server error")
	}
	return res, nil
}
func VerifyPost(postId primitive.ObjectID) error {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return errors.New("server error")
	}
	_, err = handler.Update("posts", bson.M{"_id": postId}, bson.M{"$set": bson.M{"status": "Confirmed"}})
	if err != nil {
		return errors.New("failed to verify post")
	}
	return nil
}
func RejectPost(postId primitive.ObjectID) error {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return errors.New("server error")
	}
	_, err = handler.Update("posts", bson.M{"_id": postId}, bson.M{"$set": bson.M{"status": "Rejected"}})
	if err != nil {
		return errors.New("failed to verify post")
	}
	return nil
}
