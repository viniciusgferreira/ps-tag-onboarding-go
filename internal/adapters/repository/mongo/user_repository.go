package mongo

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

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user with same first and last name already exists")
)

func (ur *UserMongoRepository) FindById(ctx *gin.Context, id string) (*models.User, error) {
	filter := bson.D{{"_id", id}}
	var user *models.User
	err := ur.db.Database(ur.dbName).Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		} else {
			return nil, err
		}
	}
	return user, nil
}
func (ur *UserMongoRepository) Save(ctx *gin.Context, u models.User) (*models.User, error) {
	filter := bson.D{{"_id", u.ID}}
	update := bson.D{{"$set", bson.D{{"firstname", u.FirstName}, {"lastname", u.LastName}, {"email", u.Email}, {"age", u.Age}}}}
	opts := options.Update().SetUpsert(true)
	_, err := ur.db.Database(ur.dbName).Collection("users").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (ur *UserMongoRepository) Update(ctx *gin.Context, u models.User) (*models.User, error) {
	filter := bson.D{{"_id", u.ID}}
	update := bson.D{{"$set", bson.D{{"firstname", u.FirstName}, {"lastname", u.LastName}, {"email", u.Email}, {"age", u.Age}}}}
	_, err := ur.db.Database(ur.dbName).Collection("users").UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
func (ur *UserMongoRepository) ExistsByFirstNameAndLastName(ctx *gin.Context, firstName, lastName string) (bool, error) {
	filter := bson.D{{"firstname", firstName}, {"lastname", lastName}}
	count, err := ur.db.Database(ur.dbName).Collection("users").CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, ErrUserAlreadyExists
	}
	return false, nil
}

func New(db *mongo.Client, dbName string) *UserMongoRepository {
	return &UserMongoRepository{db: db, dbName: dbName}
}

func Connect(db config.DB) *mongo.Client {
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
	return conn
}
