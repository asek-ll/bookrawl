package authors

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Author struct {
	Id   int `bson:"_id"`
	Name string
}

type Store struct {
	Collection *mongo.Collection
}

func (store *Store) Upsert(author *Author) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{Key: "_id", Value: author.Id}}
	update := bson.D{{Key: "$set", Value: author}}

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

func isNameMatchedStrict(name string, author *Author) bool {
	if author == nil {
		return false
	}

	authorsWords := make(map[string]bool)
	for _, s := range strings.Split(author.Name, " ") {
		authorsWords[strings.ToLower(s)] = true
	}

	for _, s := range strings.Split(name, " ") {
		if _, e := authorsWords[strings.ToLower(s)]; !e {
			return false
		}
	}

	return true
}

func (store *Store) FindByName(name string) (*Author, error) {
	filter := bson.D{{"$text", bson.D{{"$search", name}}}}
	sort := bson.D{{"score", bson.D{{"$meta", "textScore"}}}}

	opts := options.Find().SetSort(sort).SetLimit(1)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := store.Collection.Find(ctx, filter, opts)

	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var models []Author
	err = cursor.All(ctx, &models)

	if err != nil {
		return nil, err
	}

	if len(models) > 0 {
		author := &models[0]
		if isNameMatchedStrict(name, author) {
			return author, nil
		}
	}

	return nil, nil
}
