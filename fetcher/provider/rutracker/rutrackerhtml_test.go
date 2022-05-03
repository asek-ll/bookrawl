package rutracker

import (
	"fmt"
	"testing"
)

func TestFetch(t *testing.T) {
	s := &RutrackerHtmlScrapper{
		htmlClient: HtmlClientImpl{},
		apiClient:  &RutrackerApiClient{},
	}

	books, err := s.GetBookPage(2388, 0, 10)

	if err != nil {
		t.Error("Error on fetch", err)
	}

	if len(books) != 10 {
		t.Fatal("Too few books!")
	}

	for _, book := range books {
		fmt.Println(book)
	}
}
