package tasks

import (
	"bookrawl/app/abooks"
	"bookrawl/app/dao/authors"
	"bookrawl/fetcher/provider/abookclub"
	"bookrawl/fetcher/provider/rutracker"
	"bookrawl/fetcher/tasks"

	"go.mongodb.org/mongo-driver/mongo"
)

func CreateFetchBooksTask(client *mongo.Client) (func(), error) {
	tm := &tasks.TaskManager{
		RunnerManager: tasks.NewTaskRunManager(
			&authors.Store{
				Collection: client.Database("bookrawl").Collection("authors"),
			},
			&abookclub.AbookClubScrapper{},
			&rutracker.RutrackerRssScrapper{},
		),
		TaskStore: &tasks.TaskStore{
			Collection: client.Database("bookrawl").Collection("tasks"),
		},
		AbookStore: &abooks.AbookStore{
			Collection: client.Database("bookrawl").Collection("abooks"),
		},
	}

	return func() { tm.Process() }, nil
}
