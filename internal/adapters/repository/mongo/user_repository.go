package mongo

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/config"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log/slog"
)

const userCollection = "users"

type UserMongoRepository struct {
	db *mongo.Database
}

func New(db *mongo.Database) *UserMongoRepository {
	return &UserMongoRepository{db: db}
}

func (ur *UserMongoRepository) FindById(ctx *gin.Context, id string) (*models.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", oid}}
	var user *models.User
	err = ur.db.Collection(userCollection).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (ur *UserMongoRepository) Save(ctx *gin.Context, u models.User) (*models.User, error) {
	result, err := ur.db.Collection(userCollection).InsertOne(ctx, u)
	if err != nil {
		return nil, err
	}
	u.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return &u, nil
}

func (ur *UserMongoRepository) Update(ctx *gin.Context, u models.User) (*models.User, error) {
	oid, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		return nil, err
	}
	updatedUser := &models.User{}
	filter := bson.D{{"_id", oid}}
	ret := options.ReturnDocument(1)
	opts := options.FindOneAndUpdateOptions{ReturnDocument: &ret}
	update := bson.D{{"$set", bson.D{{"firstName", u.FirstName}, {"lastName", u.LastName}, {"email", u.Email}, {"age", u.Age}}}}
	err = ur.db.Collection(userCollection).FindOneAndUpdate(ctx, filter, update, &opts).Decode(updatedUser)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}
func (ur *UserMongoRepository) ExistsByFirstNameAndLastName(ctx *gin.Context, firstName, lastName string) (bool, error) {
	filter := bson.D{{"firstName", firstName}, {"lastName", lastName}}
	count, err := ur.db.Collection(userCollection).CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func Connect(db config.DB) *mongo.Database {
	slog.Info("Connecting to mongodb database")
	opts := options.Client().ApplyURI(db.Uri).SetAuth(
		options.Credential{
			Username: db.User,
			Password: db.Password,
		})
	conn, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		slog.Error("connecting to database", "error", err)
		panic(err)
	}
	if err := conn.Ping(context.TODO(), readpref.Primary()); err != nil {
		slog.Error("pinging database", "error", err)
		panic(err)
	}
	slog.Info("Database Connected")
	return conn.Database(db.Name)
}
