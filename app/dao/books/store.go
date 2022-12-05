package books

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	Collection *mongo.Collection
}

func (store *Store) Create(book *Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := store.Collection.InsertOne(ctx, book)

	if err != nil {
		return err
	}

	book.Id = result.InsertedID.(primitive.ObjectID)

	return nil
}

func (store *Store) Update(book *Book) error {
	filter := bson.D{
		{Key: "_id", Value: book.Id},
	}

	update := bson.D{{Key: "$set", Value: book}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := store.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (store *Store) Upsert(book *Book) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.D{
		{Key: "_id", Value: book.Id},
	}

	update := bson.D{{Key: "$set", Value: book}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := store.Collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (store *Store) FindByFantLabId(fantlabId int) (*Book, error) {
	filter := bson.D{
		{Key: "fantlabId", Value: fantlabId},
	}

	opts := options.Find().SetLimit(1)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := store.Collection.Find(ctx, filter, opts)

	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var models []Book
	err = cursor.All(ctx, &models)

	if err != nil {
		return nil, err
	}

	if len(models) > 0 {
		state := &models[0]
		return state, nil
	}

	return nil, nil
}
