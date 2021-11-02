package tasks

import (
	"bookrawl/app/abooks"
	"bookrawl/fetcher/provider/abookclub"
	"bookrawl/fetcher/provider/rutracker"
	"bookrawl/fetcher/tasks"

	"go.mongodb.org/mongo-driver/mongo"
)

func fetchBooks(client *mongo.Client) error {
	tm := &tasks.TaskManager{
		RunnerManager: tasks.NewTaskRunManager(
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

	return tm.Process()
}