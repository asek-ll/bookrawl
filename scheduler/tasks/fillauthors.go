package tasks

import (
	"bookrawl/app/dao/abooks"
	"bookrawl/app/dao/authors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func FillAuthors(client *mongo.Client) error {
	authorStore := &authors.Store{
		Collection: client.Database("bookrawl").Collection("authors"),
	}

	abookStore := &abooks.AbookStore{
		Collection: client.Database("bookrawl").Collection("abooks"),
	}

	var filter *abooks.FindBooksFilter = nil

	for true {
		page, err := abookStore.Find(filter, 100)
		if err != nil {
			return err
		}
		var lastDate *time.Time = nil
		for _, book := range page.Books {
			if book.AuthorId == nil {
				author, err := authorStore.FindByName(book.Author)
				if err != nil {
					return err
				}
				if author != nil {
					log.Println("find author for for book", book.Author, author, author.Id)
					book.AuthorId = []int{author.Id}
					err = abookStore.Upsert(book)
					if err != nil {
						return err
					}
				} else {
					log.Println("Can't find author for for book", book.Author)
				}
			}
			lastDate = &book.Date
		}
		if lastDate == nil {
			break
		}
		filter = &abooks.FindBooksFilter{
			BeforeDate: lastDate,
		}
	}
	return nil
}
