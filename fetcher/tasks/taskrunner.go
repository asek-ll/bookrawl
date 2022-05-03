package tasks

import (
	"bookrawl/app/dao/abooks"
	"time"
)

type TaskRunner interface {
	GetType() string
	Fetch(task SyncTask, start time.Time) ([]abooks.ABook, error)
}
