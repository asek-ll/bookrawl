package dao

import (
	"bookrawl/app/dao/abooks"
	"bookrawl/app/dao/authors"

	"go.mongodb.org/mongo-driver/mongo"
)

type DaoHolder struct {
	client *mongo.Client
	bookStore *abooks.AbookStore
	authorStore *authors.Store
}

func NewDaoHolder(client *mongo.Client) *DaoHolder {
	return &DaoHolder{
		client: client,
		bookStore: &abooks.AbookStore{
			Collection: client.Database("bookrawl").Collection("abooks"),
		},
		authorStore: &authors.Store{
			Collection: client.Database("bookrawl").Collection("authors"),
		},
	}
}


func (dh *DaoHolder) GetBookStore() *abooks.AbookStore {
	return dh.bookStore
}

func (dh *DaoHolder) GetAuthorsStore() *authors.Store {
	return dh.authorStore
}
