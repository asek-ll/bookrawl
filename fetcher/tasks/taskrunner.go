package tasks

import (
	"bookrawl/app/dao/abooks"
)

type TaskRunner interface {
	GetType() string
	Fetch(params TaskParams) ([]abooks.ABook, error)
}
