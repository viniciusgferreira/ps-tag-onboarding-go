package config

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log/slog"
	"time"
)

func Connect(ctx context.Context, db DB) (*mongo.Database, error) {
	slog.Info("Connecting to mongodb database")
	opts := options.Client().ApplyURI(db.Uri).SetAuth(
		options.Credential{
			Username: db.User,
			Password: db.Password,
		})
	conn, err := mongo.Connect(ctx, opts)
	if err != nil {
		slog.Error("connecting to database", "error", err)
		return nil, err
	}
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := conn.Ping(pingCtx, readpref.Primary()); err != nil {
		slog.Error("pinging database", "error", err)
		return nil, err
	}
	slog.Info("Database Connected")
	return conn.Database(db.Name), nil
}
