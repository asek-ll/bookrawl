package tasks

import (
    "fmt"
    "encoding/xml"
    "net/http"
    "io/ioutil"
    "strings"
    "time"
)

const (
    rutrackerTimeFormat string = "2006-01-02T15:04:05"
)

type EntryLink struct{
    Href string  `xml:"href,attr"`
}

type AtomEntry struct {
    Id string `xml:"id"`
    Link EntryLink `xml:"link"`
    Updated string `xml:"updated"`
    Title string `xml:"title"`
}

type AtomFeed struct {
    Title string `xml:"title,attr"`
    Entries []AtomEntry `xml:"entry"`
}

type RutrackerRssScrapper struct {
}

func (s *RutrackerRssScrapper) GetType() string {
    return "rutracker"
}

func (s *RutrackerRssScrapper) Fetch(params TaskParams) ([]ABook, error) {
    var forumId string
    exists := false
    if params != nil {
        forumId, exists = params["forumId"]
    }
    if !exists {
        return nil, fmt.Errorf("No forum id present")
    } 
    url := fmt.Sprintf("http://feed.rutracker.cc/atom/f/%s.atom", forumId)
    response, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()
    body, err := ioutil.ReadAll(response.Body)
    var feed AtomFeed
    err = xml.Unmarshal(body, &feed)
    if err != nil {
        return nil, err
    }


    books := make([]ABook, len(feed.Entries))

    for i, entry := range feed.Entries {

        id := strings.Split(entry.Link.Href, "=")[1]
        rawTitleParts := strings.SplitN(entry.Title, " - ", 2)
        if len(rawTitleParts) < 2 {
            rawTitleParts = strings.SplitN(entry.Title, " – ", 2)
        }
        var title, author string
        if len(rawTitleParts) >= 2 {
            author = rawTitleParts[0]
            title = rawTitleParts[1]
        } else {
            title = rawTitleParts[0]
        }

        updated, _ := time.Parse(rutrackerTimeFormat, entry.Updated[:len(rutrackerTimeFormat)])

        book := ABook{
            Id: fmt.Sprintf("rutracker-%s", id),
            RawTitle: entry.Title,
            Title: title,
            Author: author,
            Date: updated,
            Link: entry.Link.Href,
        }

        books[i] = book
    }

    return books, nil
}
