package abooks

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AbooksPage struct {
	Books    []ABook
	HasNext  bool
	PageSize int
}

type AbookStore struct {
	Collection *mongo.Collection
}

func (as *AbookStore) InsertBooks(books []ABook) error {

	models := make([]interface{}, len(books))
	for i, book := range books {
		models[i] = book
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := as.Collection.InsertMany(ctx, models)

	if err != nil {
		return err
	}

	if len(res.InsertedIDs) != len(books) {
		return fmt.Errorf("Can't insert books %v", books)
	}

	return nil
}

func (as *AbookStore) Upsert(book ABook) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{Key: "id", Value: book.Id}}
	update := bson.D{{Key: "$set", Value: book}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := as.Collection.UpdateOne(ctx, filter, update, opts)

	if err != nil {
		return err
	}

	if res.MatchedCount == 0 && res.UpsertedCount == 0 {
		return fmt.Errorf("Can't insert book %v", book)
	}

	return nil
}

func (as *AbookStore) UpsertMany(books []ABook) error {
	models := make([]mongo.WriteModel, len(books))

	for i, book := range books {
		filter := bson.D{{Key: "id", Value: book.Id}}
		update := bson.D{{Key: "$set", Value: book}}
		models[i] = mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
	}

	opts := options.BulkWrite().SetOrdered(false)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := as.Collection.BulkWrite(ctx, models, opts)
	return err
}

func (as *AbookStore) GetById(id string) (*ABook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var model ABook
	err := as.Collection.FindOne(ctx, bson.M{"id": id}).Decode(&model)

	if err != nil {
		return nil, err
	}

	return &model, nil

}

func (as *AbookStore) Find(filter *FindBooksFilter, pageSize int) (*AbooksPage, error) {
	opts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}}).SetLimit(int64(pageSize + 1))
	queryFilter := bson.D{}

	if filter != nil {
		if filter.AfterDate != nil {
			queryFilter = append(queryFilter, bson.E{
				Key: "date",
				Value: bson.D{{
					Key:   "$gte",
					Value: *filter.AfterDate},
				},
			})
		}

		if filter.BeforeDate != nil {
			queryFilter = append(queryFilter, bson.E{
				Key: "date",
				Value: bson.D{{
					Key:   "$lt",
					Value: *filter.BeforeDate},
				},
			})
		}
		if filter.AuthorId != nil {
			queryFilter = append(queryFilter, bson.E{
				Key:   "authorid",
				Value: *filter.AuthorId,
			})
		} else if filter.NoAuthor != nil && *filter.NoAuthor == true {
			queryFilter = append(queryFilter, bson.E{
				Key:   "authorid",
				Value: nil,
			})
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := as.Collection.Find(ctx, queryFilter, opts)

	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var models []ABook
	err = cursor.All(ctx, &models)

	if err != nil {
		return nil, err
	}

	hasNext := false
	if len(models) > pageSize {
		hasNext = true
		models = models[:pageSize]
	}

	result := &AbooksPage{models, hasNext, pageSize}

	return result, nil
}
