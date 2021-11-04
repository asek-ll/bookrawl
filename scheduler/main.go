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
	_, err = s.Cron("0 1 * * *").Do(fetchAuthorsTask)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Register 'fetchAuthorsTask' task")

	fetchBooksTask, _ := tasks.CreateFetchBooksTask(client)
	_, err = s.Cron("0 3 * * *").Do(fetchBooksTask)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Register 'fetchBooksTask' task")

	// err = tasks.FillAuthors(client)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	s.StartBlocking()

}
