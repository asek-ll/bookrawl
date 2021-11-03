package tasks

import (
	"bookrawl/app/dao/authors"
	"bookrawl/app/fantlab"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

func CreateFetchAuthorsTask(client *mongo.Client) (func(), error) {
	authorStore := &authors.Store{
		Collection: client.Database("bookrawl").Collection("authors"),
	}

	return func() { fetchAuthors(&fantlab.Api{}, authorStore) }, nil
}

func fetchAuthors(api *fantlab.Api, store *authors.Store) error {
	log.Printf("Start 'fetchAuthorsTask' task")

	resp, err := api.GetAllAuthors()

	if err != nil {
		return err
	}

	buf := make([]*authors.Author, 100)
	bufLen := 0

	log.Printf("Recieved %d authors", len(resp.List))

	for _, author := range resp.List {
		authorDto := &authors.Author{
			Id:   author.Id,
			Name: author.Name,
		}
		buf[bufLen] = authorDto
		bufLen += 1

		if bufLen == len(buf) {
			err = store.UpsertMany(buf)
			if err != nil {
				return err
			}
			bufLen = 0
		}
	}

	if bufLen > 0 {
		err = store.UpsertMany(buf[0:bufLen])
		if err != nil {
			return err
		}
	}

	log.Printf("Complete 'fetchAuthorsTask' task")

	return nil
}
