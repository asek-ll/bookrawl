package tasks

import (
	"bookrawl/app/abooks"
)

type TaskRunner interface {
	GetType() string
	Fetch(params TaskParams) ([]abooks.ABook, error)
}
