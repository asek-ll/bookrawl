package abooks 

import (
	"time"
)

type FindBooksFilter struct {
	AfterDate  *time.Time
	BeforeDate *time.Time
	AuthorId *int
}

type BooksFilterBuilder struct {
	filter *FindBooksFilter
}

func NewBooksFilterBuilder() *BooksFilterBuilder {
	return &BooksFilterBuilder{
		filter: &FindBooksFilter{
			AfterDate: nil,
			BeforeDate: nil,
			AuthorId: nil,
		},
	}
}

func (fb *BooksFilterBuilder) SetAfterDate(t *time.Time) *BooksFilterBuilder {
	fb.filter.AfterDate = t
	return fb
}

func (fb *BooksFilterBuilder) SetBeforeDate(t *time.Time) *BooksFilterBuilder {
	fb.filter.BeforeDate = t
	return fb
}

func (fb *BooksFilterBuilder) SetAuthorId(authorId *int) *BooksFilterBuilder {
	fb.filter.AuthorId = authorId
	return fb
}

func (fb *BooksFilterBuilder) Build() *FindBooksFilter {
	return fb.filter
}
