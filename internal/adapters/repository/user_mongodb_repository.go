package repository

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
)

const userCollection = "users"

type UserMongoRepository struct {
	db *mongo.Database
}

func NewUserRepo(db *mongo.Database) *UserMongoRepository {
	return &UserMongoRepository{db: db}
}

func (ur *UserMongoRepository) FindById(ctx *gin.Context, id string) (*model.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		slog.Error("converting user id from request to object id.", "error", err)
		return nil, nil
	}
	filter := bson.D{{"_id", oid}}
	var user *model.User
	err = ur.db.Collection(userCollection).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		slog.Error("failed to decode FindOne result", "error", err)
		return nil, err
	}
	return user, nil
}
func (ur *UserMongoRepository) Save(ctx *gin.Context, u model.User) (*model.User, error) {
	result, err := ur.db.Collection(userCollection).InsertOne(ctx, u)
	if err != nil {
		slog.Error("failed to insert user", "error", err)
		return nil, err
	}
	u.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return &u, nil
}

func (ur *UserMongoRepository) Update(ctx *gin.Context, u model.User) (*model.User, error) {
	oid, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		return nil, err
	}
	updatedUser := &model.User{}
	filter := bson.D{{"_id", oid}}
	ret := options.ReturnDocument(1)
	opts := options.FindOneAndUpdateOptions{ReturnDocument: &ret}
	update := bson.D{{"$set", bson.D{{"firstName", u.FirstName}, {"lastName", u.LastName}, {"email", u.Email}, {"age", u.Age}}}}
	err = ur.db.Collection(userCollection).FindOneAndUpdate(ctx, filter, update, &opts).Decode(updatedUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		slog.Error("failed to decode FindOneAndUpdate result", "error", err)
		return nil, err
	}
	return updatedUser, nil
}
func (ur *UserMongoRepository) ExistsByFirstNameAndLastName(ctx *gin.Context, u model.User) (bool, error) {
	var oid primitive.ObjectID
	var err error
	if len(u.ID) != 0 {
		oid, err = primitive.ObjectIDFromHex(u.ID)
		if err != nil {
			return false, err
		}
	}
	filter := bson.D{
		{"firstName", u.FirstName},
		{"lastName", u.LastName},
		{"_id", bson.D{{"$ne", oid}}},
	}
	var user *model.User
	err = ur.db.Collection(userCollection).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		slog.Error("failed to decode FindOne result", "error", err)
		return false, err
	}
	if user != nil {
		return true, nil
	}
	return false, nil
}
