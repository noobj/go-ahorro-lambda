package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	. "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AbstractRepository struct {
	BaseRepository
	IRepository
}

func (repo AbstractRepository) InsertOne(doc interface{}) {
	_, err := repo.Collection.InsertOne(context.TODO(), doc)
	if err != nil {
		panic(err)
	}
}

func (repo AbstractRepository) Disconnect() func() {
	return func() {
		if err := repo.Client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}
}

func (repo AbstractRepository) Aggregate(stages interface{}) []bson.M {
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

func (repo AbstractRepository) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *SingleResult {
	return repo.Collection.FindOne(ctx, filter, opts...)
}
