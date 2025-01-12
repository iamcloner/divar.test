package userRepositories

import (
	"divar.ir/internal/mongodb"
	"divar.ir/schema"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func GetMyPosts(userId primitive.ObjectID) (posts []bson.M, err error) {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return nil, errors.New("internal server error")
	}
	result, err := handler.FindMany("posts", bson.M{"userId": userId}, bson.M{})
	if err != nil {
		return nil, errors.New("nothing found")
	}
	return result, nil
}
func AddPost(post schema.Posts) error {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return errors.New("internal server error")
	}
	_, err = handler.Insert("posts", post)
	if err != nil {
		return errors.New("nothing found")
	}
	return nil
}
func UpdatePost(post schema.Posts) error {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return errors.New("internal server error")
	}
	var oldPost schema.Posts
	err = handler.FindOne("posts", bson.M{"_id": post.ID, "userId": post.UserId, "status": bson.M{"$in": []string{"Confirmed", "Pending"}}}, bson.M{"createdAt": 1}).Decode(&oldPost)

	if err != nil {
		return errors.New("nothing found")
	}
	post.CreatedAt = oldPost.CreatedAt
	_, err = handler.Update("posts", bson.M{"_id": post.ID}, bson.M{"$set": post})
	if err != nil {
		return errors.New("nothing found")
	}
	return nil
}
func DeletePost(userId primitive.ObjectID, postId primitive.ObjectID) error {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return errors.New("internal server error")
	}
	var oldPost schema.Posts
	err = handler.FindOne("posts", bson.M{"_id": postId, "userId": userId, "status": bson.M{"$in": []string{"Confirmed", "Pending"}}}, bson.M{}).Decode(&oldPost)
	if err != nil {
		return errors.New("nothing found")
	}
	_, err = handler.Update("posts", bson.M{"_id": postId}, bson.M{"$set": bson.M{"status": "deleted", "updatedAt": time.Now(), "lastAction": "deleted by user"}})
	if err != nil {
		return errors.New("nothing found")
	}
	return nil
}
