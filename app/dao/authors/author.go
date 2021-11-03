package authors

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Author struct {
	Id   int
	Name string
}

type Store struct {
	Collection *mongo.Collection
}

func (store *Store) Upsert(author *Author) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{Key: "_id", Value: author.Id}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{"_id", author.Id},
		{"name", author.Name},
	}}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := store.Collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (store *Store) UpsertMany(authors []*Author) error {
	models := make([]mongo.WriteModel, len(authors))

	for i, author := range authors {
		filter := bson.D{{Key: "_id", Value: author.Id}}
		update := bson.D{{Key: "$set", Value: bson.D{
			{"_id", author.Id},
			{"name", author.Name},
		}}}
		models[i] = mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
	}

	opts := options.BulkWrite().SetOrdered(false)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := store.Collection.BulkWrite(ctx, models, opts)
	return err
}

func (store *Store) FindByName(name string) (*Author, error) {
	return nil, nil
}
