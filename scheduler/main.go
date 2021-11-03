package main

import (
	"bookrawl/scheduler/tasks"
	"bookrawl/scheduler/utils"
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	s := gocron.NewScheduler(time.UTC)

	client, err := utils.CreateMongoClient()

	if err != nil {
		log.Fatal(err)
	}

	fetchAuthorsTask, _ := tasks.CreateFetchAuthorsTask(client)
	s.Cron("0 7 * * *").Do(fetchAuthorsTask)
	log.Printf("Register 'fetchAuthorsTask' task")

	fetchBooksTask, _ := tasks.CreateFetchBooksTask(client)
	s.Cron("0 5 * * *").Do(fetchBooksTask)
	log.Printf("Register 'fetchBooksTask' task")

	s.StartBlocking()
}
