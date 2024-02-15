package repository

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/config"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log/slog"
)

type UserMongoRepository struct {
	db     *mongo.Client
	dbName string
}

var UserNotFound = errors.New("user not found")

func (ur *UserMongoRepository) FindById(ctx *gin.Context, id string) (*models.User, error) {
	filter := bson.D{{"_id", id}}
	var user *models.User
	err := ur.db.Database(ur.dbName).Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, UserNotFound
		} else {
			return nil, err
		}
	}
	return user, nil
}
func (ur *UserMongoRepository) FindAll(ctx *gin.Context) ([]*models.User, error) {
	return nil, nil
}
func (ur *UserMongoRepository) Save(ctx *gin.Context, u *models.User) (*models.User, error) {
	filter := bson.D{{"_id", u.ID}}
	update := bson.D{{"$set", bson.D{{"firstname", u.FirstName}, {"lastname", u.LastName}, {"email", u.Email}, {"age", u.Age}}}}
	opts := options.Update().SetUpsert(true)

	_, err := ur.db.Database(ur.dbName).Collection("users").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func New(db *mongo.Client) *UserMongoRepository {
	return &UserMongoRepository{db: db, dbName: "onboardingdb"}
}

func Credentials(db *config.DB) options.Credential {
	return options.Credential{
		Username: db.User,
		Password: db.Password,
	}
}

func Connect(db *config.DB) *mongo.Client {
	slog.Info("Connecting to mongodb database")
	opts := options.Client().ApplyURI(db.Uri).SetAuth(Credentials(db))
	conn, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		slog.Error("connecting to database", "error", err)
		panic("database error")
	}
	if err := conn.Ping(context.TODO(), readpref.Primary()); err != nil {
		slog.Error("pinging database", "error", err)
		panic("database error")
	}
	slog.Info("Database Connected")
	return conn
}
