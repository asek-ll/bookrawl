package authors

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Author struct {
	Id      int
	Name    string
	Tracked bool
}

type Store struct {
	Collection *mongo.Collection
}

func (store *Store) Upsert(author *Author) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{Key: "id", Value: author.Id}}
	update := bson.D{{Key: "$set", Value: author}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := as.Collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (store *Store) FindByName(name String) (*Author, error) {
	return nil, nil
}
