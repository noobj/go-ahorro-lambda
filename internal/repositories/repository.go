package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:generate mockgen -source=repository.go -destination mocks/mock_repository.go
type IRepository interface {
	Aggregate(interface{}) []bson.M
	// TODO: use generic for Entry type, not bson
	InsertOne(context.Context, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	InsertMany(context.Context, []interface{}, ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	Disconnect() func()
	FindOne(context.Context, interface{}, ...*options.FindOneOptions) *mongo.SingleResult
	DeleteMany(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}

type BaseRepository struct {
	Collection *mongo.Collection
}

func (repo BaseRepository) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return repo.Collection.InsertOne(context.TODO(), document, opts...)
}

func (repo BaseRepository) Disconnect() func() {
	return func() {
		if err := repo.Collection.Database().Client().Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}
}

func (repo BaseRepository) Aggregate(stages interface{}) []bson.M {
	cursor, err := repo.Collection.Aggregate(context.TODO(), stages)
	if err != nil {
		panic(err)
	}

	var results []bson.M

	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results
}

func (repo BaseRepository) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return repo.Collection.FindOne(ctx, filter, opts...)
}

func (repo BaseRepository) DeleteMany(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return repo.Collection.DeleteMany(ctx, filter, opts...)
}

func (repo BaseRepository) InsertMany(ctx context.Context, documents []interface{},
	opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return repo.Collection.InsertMany(ctx, documents, opts...)
}
