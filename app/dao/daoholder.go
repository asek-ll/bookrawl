package dao

import (
	"bookrawl/app/dao/abooks"
	"bookrawl/app/dao/authors"
	"bookrawl/app/dao/users"

	"go.mongodb.org/mongo-driver/mongo"
)

type DaoHolder struct {
	client      *mongo.Client
	bookStore   *abooks.AbookStore
	authorStore *authors.Store
	userStore   *users.Store
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
		userStore: &users.Store{
			Collection: client.Database("bookrawl").Collection("users"),
		},
	}
}

func (dh *DaoHolder) GetBookStore() *abooks.AbookStore {
	return dh.bookStore
}

func (dh *DaoHolder) GetAuthorsStore() *authors.Store {
	return dh.authorStore
}

func (dh *DaoHolder) GetUsersStore() *users.Store {
	return dh.userStore
}
