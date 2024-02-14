package repository

import (
	"context"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/config"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
)

type UserMongoRepository struct {
	db *mongo.Client
}

func (ur *UserMongoRepository) FindById(id string) (models.User, error) {
	return models.User{}, nil
}
func (ur *UserMongoRepository) FindAll() ([]models.User, error) {
	return []models.User{}, nil
}
func (ur *UserMongoRepository) Save(u models.User) (models.User, error) {
	return models.User{}, nil
}

func New(db *mongo.Client) *UserMongoRepository {
	return &UserMongoRepository{db: db}
}

func Credentials(db *config.DB) options.Credential {
	return options.Credential{
		Username: db.User,
		Password: db.Password,
	}
}

func Connect(db *config.DB) (*mongo.Client, error) {
	slog.Info("Connecting to mongodb database")
	opts := options.Client().ApplyURI(db.Uri).SetAuth(Credentials(db))
	ctx := context.TODO()
	conn, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	return conn, nil

}
