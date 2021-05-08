package tasks

import (
	"time"
)

type ABook struct {
	Id          string
	RawTitle    string
	Title       string
	Author      string
	Artist      string
	Year        int
	Date        time.Time
	Link        string
	Description string
}
