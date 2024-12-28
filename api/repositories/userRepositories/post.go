package userRepositories

import (
	"divar.ir/api/repositories"
	"divar.ir/internal/mongodb"
	"divar.ir/schema"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
func AddPost(userId primitive.ObjectID, post schema.Posts) error {
	handler, err := mongodb.GetMongoDBHandler()
	if err != nil {
		return errors.New("internal server error")
	}
	post.UserId = userId
	_, err = handler.Insert("posts", post)
	if err != nil {
		return errors.New("nothing found")
	}
	return nil
}
func CheckPostRequirements(post schema.Posts) error {
	if post.CategoryCode == 0 {
		return errors.New("category code is required")
	}
	result, err := repositories.GetCategories(post.CategoryCode)
	if err != nil {
		return errors.New("invalid Category code")
	}

	println(result)
	return nil
}
