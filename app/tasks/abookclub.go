package tasks

import (
    "github.com/PuerkitoBio/goquery"
    "net/http"
    "strings"
    "io"
    "time"
    "fmt"
)

const (
    timeFormat string = "02.01.2006 15:04"
)

type AbookClubScrapper struct {
}

func (s *AbookClubScrapper) Fetch(params TaskParams) ([]ABook, error) {
    response, err := http.Get("http://abook-club.ru/new_abooks/")
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()

    return parseBody(response.Body)
}

func (s *AbookClubScrapper) GetType() string {
    return "abook-club"
}

func parseBody(reader io.Reader) ([]ABook, error) {

    // Create a goquery document from the HTTP response
    document, err := goquery.NewDocumentFromReader(reader)
    if err != nil {
        return nil, err
    }

    books := []ABook{}

    document.Find("div.entry").Each(func(index int, el *goquery.Selection) {
        a := el.Find("div.entry_header_full a")
        rawTitle := a.Text()
        rawTitleParts := strings.Split(rawTitle, " - ")

        link, _ := a.Attr("href")
        id := strings.Split(link, "=")[1]
        authorAndTime := el.Find("div.entry_time").Text()
        authorAndTimeParts := strings.Split(authorAndTime, ", ")
        description, _ := el.Find("div.entry_content").Html()

        datetime := strings.Split(authorAndTimeParts[1], "\n")[0]

        created, _ := time.Parse(timeFormat, datetime)

        book := ABook{
            Id: fmt.Sprintf("abook-club-%s", id),
            RawTitle: rawTitle,
            Title: rawTitleParts[1],
            Author: rawTitleParts[0],
            Date: created,
            Link: link,
            Description: description,
        }

        //fmt.Println(book.Id, book.Title, book.Date)

        books = append(books, book)

    })

    if len(books) == 0 {
        return nil, fmt.Errorf("No books founded")
    }

    return books, nil
}



