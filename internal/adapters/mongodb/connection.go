package mongodb

import (
	"context"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/adapters/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log/slog"
)

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