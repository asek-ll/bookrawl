package userbookstates

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

func (store *Store) Upsert(state *State) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.D{
		{Key: "bookId", Value: state.BookId},
		{Key: "userId", Value: state.UserId},
	}

	update := bson.D{{Key: "$set", Value: state}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := store.Collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (store *Store) FindByBookIdAndUserId(bookId, userId primitive.ObjectID) (*State, error) {
	filter := bson.D{
		{Key: "bookId", Value: bookId},
		{Key: "userId", Value: userId},
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
	var models []State
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
