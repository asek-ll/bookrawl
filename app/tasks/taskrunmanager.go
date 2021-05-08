package tasks

import (
	//"fmt"
	"errors"
	"log"
	"time"
)

type TaskRunManager struct {
	RunnerByType map[string]TaskRunner
}

func NewTaskRunManager(runners ...TaskRunner) *TaskRunManager {
	rm := TaskRunManager{RunnerByType: make(map[string]TaskRunner)}

	for _, runner := range runners {
		rm.RunnerByType[runner.GetType()] = runner
	}

	return &rm
}

func (tm *TaskRunManager) run(task SyncTask, start time.Time) ([]ABook, error) {
	runner, ex := tm.RunnerByType[task.Type]
	if !ex {
		return nil, errors.New("No runner")
	}

	log.Printf("Process '%s' runner", task.Type)

	books, err := runner.Fetch(task.Params)
	if err != nil {
		log.Printf("Error in runner - %s", err.Error())
		return nil, err
	}

	log.Printf("Fetched books count %d", len(books))

	filtered := []ABook{}

	for _, book := range books {
		if book.Date.After(task.LastRun) && !book.Date.After(start) {
			filtered = append(filtered, book)
		}
	}

	log.Printf("Filtered books count %d", len(filtered))

	return filtered, nil
}

func (tm *TaskRunManager) RunTask(task SyncTask) (SyncTask, []ABook) {
	now := time.Now()
	books, err := tm.run(task, now)

	updatedTask := task

	if err != nil {
		updatedTask.State = "error"
		updatedTask.ErrorMsg = err.Error()
		return updatedTask, []ABook{}
	}

	updatedTask.LastRun = now
	updatedTask.State = "success"
	updatedTask.ErrorMsg = ""
	return updatedTask, books
}
