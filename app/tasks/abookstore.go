package tasks

import (
    "fmt"
    "time"
    "context"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
)

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
