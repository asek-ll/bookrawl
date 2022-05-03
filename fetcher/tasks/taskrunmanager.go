package tasks

import (
	//"fmt"
	"errors"
	"log"
	"time"

	"bookrawl/app/dao/abooks"
	"bookrawl/app/dao/authors"
)

type TaskRunManager struct {
	RunnerByType map[string]TaskRunner
	AuthorsStore *authors.Store
}

func NewTaskRunManager(authorsStore *authors.Store, runners ...TaskRunner) *TaskRunManager {
	rm := TaskRunManager{
		RunnerByType: make(map[string]TaskRunner),
		AuthorsStore: authorsStore,
	}

	for _, runner := range runners {
		rm.RunnerByType[runner.GetType()] = runner
	}

	return &rm
}

func (tm *TaskRunManager) run(task SyncTask, start time.Time) ([]abooks.ABook, error) {
	runner, ex := tm.RunnerByType[task.Type]
	if !ex {
		return nil, errors.New("No runner")
	}

	log.Printf("Process '%s' runner", task.Type)

	books, err := runner.Fetch(task, start)
	if err != nil {
		log.Printf("Error in runner - %s", err.Error())
		return nil, err
	}

	log.Printf("Fetched books count %d", len(books))

	filtered := []abooks.ABook{}

	for _, book := range books {
		author, err := tm.AuthorsStore.FindByName(book.Author)
		if err != nil {
			log.Println("Error fetch author for book", book.RawTitle, err)
		} else if author == nil {
			log.Printf("No author [%s] founded for book: %s", book.Author, book.RawTitle)
		} else {
			book.AuthorId = []int{author.Id}
		}
		filtered = append(filtered, book)
	}

	log.Printf("Filtered books count %d", len(filtered))

	return filtered, nil
}

func (tm *TaskRunManager) RunTask(task SyncTask) (SyncTask, []abooks.ABook) {
	now := time.Now()
	books, err := tm.run(task, now)

	updatedTask := task

	if err != nil {
		updatedTask.State = "error"
		updatedTask.ErrorMsg = err.Error()
		return updatedTask, []abooks.ABook{}
	}

	updatedTask.LastRun = now
	updatedTask.State = "success"
	updatedTask.ErrorMsg = ""
	return updatedTask, books
}
