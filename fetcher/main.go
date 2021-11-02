package main

import (
	"bookrawl/app/abooks"
	"bookrawl/fetcher/provider/abookclub"
	"bookrawl/fetcher/provider/rutracker"
	"bookrawl/fetcher/tasks"
	"os"

	"context"
	"flag"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	connUri := flag.String("mongodbURI", "", "mongodb connection uri")
	flag.Parse()

	if *connUri == "" {
		*connUri = os.Getenv("MONGODB_URI")
	}

	if *connUri == "" {

		log.Fatal("mongodbURI required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(*connUri))
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		log.Fatal(err)
	}

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

	err = tm.Process()
	if err != nil {
		log.Fatal(err)
	}

}