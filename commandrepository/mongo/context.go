package mongo

import (
	"context"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EnsureIndexes создаёт необходимые индексы в коллекциях Event Store.
func EnsureIndexes(ctx context.Context, db *mongodriver.Database, opts Options, logger *slog.Logger) error {
	// Индексы для событий: compound (id + version) unique, correlation_token, created_at
	eventsColl := db.Collection(opts.EventsCollection)
	_, err := eventsColl.Indexes().CreateMany(ctx, []mongodriver.IndexModel{
		{
			Keys:    bson.D{{Key: "_id", Value: 1}, {Key: "version", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{Keys: bson.D{{Key: "correlation_token", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: 1}}},
	})
	if err != nil {
		return fmt.Errorf("create events indexes: %w", err)
	}

	// Индексы для snapshot'ов
	snapshotsColl := db.Collection(opts.SnapshotsCollection)
	_, err = snapshotsColl.Indexes().CreateOne(ctx, mongodriver.IndexModel{
		Keys:    bson.D{{Key: "_id", Value: 1}, {Key: "version", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("create snapshots indexes: %w", err)
	}

	logger.Info("MongoDB indexes ensured",
		"events", opts.EventsCollection,
		"snapshots", opts.SnapshotsCollection)
	return nil
}
