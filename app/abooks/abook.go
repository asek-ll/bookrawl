package abooks

import (
	"time"
)

type ABook struct {
	Id          string
	RawTitle    string
	Title       string
	Author      string
	Artists     []string
	Year        int
	Date        time.Time
	Link        string
	Description string
	Length      int
	Size        string
	Quality     string
	Props       map[string]string
	AuthorId    []int
}
