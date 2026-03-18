package mongo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FullEventSourcingRepository — стратегия Full Event Sourcing для MongoDB.
type FullEventSourcingRepository struct {
	db     *mongodriver.Database
	opts   Options
	logger *slog.Logger
}

// NewFullEventSourcingRepository создаёт MongoDB репозиторий Full Event Sourcing.
func NewFullEventSourcingRepository(db *mongodriver.Database, opts Options, logger *slog.Logger) *FullEventSourcingRepository {
	return &FullEventSourcingRepository{db: db, opts: opts, logger: logger}
}

// GetByID получает все события агрегата, отсортированные по версии.
func (r *FullEventSourcingRepository) GetByID(ctx context.Context, id uuid.UUID) ([]DomainEventDocument, error) {
	coll := r.db.Collection(r.opts.EventsCollection)
	filter := bson.M{"_id": id}
	findOpts := options.Find().SetSort(bson.D{{Key: "version", Value: 1}})

	cursor, err := coll.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, fmt.Errorf("find events: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []DomainEventDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("decode events: %w", err)
	}
	return docs, nil
}

// GetByIDAndVersion получает событие по идентификатору и версии.
func (r *FullEventSourcingRepository) GetByIDAndVersion(ctx context.Context, id uuid.UUID, version int64) (*DomainEventDocument, error) {
	coll := r.db.Collection(r.opts.EventsCollection)
	filter := bson.M{"_id": id, "version": version}

	var doc DomainEventDocument
	if err := coll.FindOne(ctx, filter).Decode(&doc); err != nil {
		return nil, fmt.Errorf("find event: %w", err)
	}
	return &doc, nil
}

// GetByCorrelationToken получает событие по маркеру корреляции.
func (r *FullEventSourcingRepository) GetByCorrelationToken(ctx context.Context, ct uuid.UUID) (*DomainEventDocument, error) {
	coll := r.db.Collection(r.opts.EventsCollection)
	filter := bson.M{"correlation_token": ct}

	var doc DomainEventDocument
	if err := coll.FindOne(ctx, filter).Decode(&doc); err != nil {
		return nil, fmt.Errorf("find event by correlation token: %w", err)
	}
	return &doc, nil
}

// SaveEvent сохраняет событие в MongoDB.
func (r *FullEventSourcingRepository) SaveEvent(ctx context.Context, doc *DomainEventDocument) error {
	coll := r.db.Collection(r.opts.EventsCollection)
	_, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("insert event: %w", err)
	}
	r.logger.Info("FullES: event saved", "aggregateId", doc.ID, "version", doc.Version)
	return nil
}

// SaveSnapshot — не используется в стратегии FullEventSourcing.
func (r *FullEventSourcingRepository) SaveSnapshot(_ context.Context, _ uuid.UUID, _ int64, _ string) error {
	return nil
}

// GetEventCount возвращает количество событий для агрегата.
func (r *FullEventSourcingRepository) GetEventCount(ctx context.Context, id uuid.UUID) (int64, error) {
	coll := r.db.Collection(r.opts.EventsCollection)
	return coll.CountDocuments(ctx, bson.M{"_id": id})
}
