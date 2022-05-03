package users

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Id              primitive.ObjectID `bson:"_id"`
	ChatId          int64              `bson:"chatId"`
	FavoriteAuthors []int              `bson:"favoriteAuthors"`
}

type Store struct {
	Collection *mongo.Collection
}

func (s *Store) Upsert(user *User) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{Key: "_id", Value: user.Id}}
	update := bson.D{{Key: "$set", Value: user}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := s.Collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (s *Store) FindByFavoriteAuthors(favoriteAuthorIds []int) ([]User, error) {
	filter := bson.D{{"favoriteAuthors", bson.D{{"$in", favoriteAuthorIds}}}}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := s.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var models []User
	err = cursor.All(ctx, &models)

	if err != nil {
		return nil, err
	}

	return models, nil
}
