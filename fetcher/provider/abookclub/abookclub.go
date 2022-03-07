package abookclub

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bookrawl/app/dao/abooks"
	"bookrawl/fetcher/tasks"

	"github.com/PuerkitoBio/goquery"
)

const (
	timeFormat string = "02.01.2006 15:04"
)

type AbookClubScrapper struct {
}

func (s *AbookClubScrapper) GetPageBooks(page int) ([]abooks.ABook, error) {
	url := fmt.Sprintf("http://abook-club.ru/new_abooks/page=%d/", page)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return parseBody(response.Body)
}

func (s *AbookClubScrapper) Fetch(params tasks.TaskParams) ([]abooks.ABook, error) {
	firstPage, err := s.GetPageBooks(1)
	if err != nil {
		return nil, err
	}

	secondPage, err2 := s.GetPageBooks(2)

	if err2 != nil {
		return firstPage, nil
	}

	return append(firstPage, secondPage...), nil
}

func (s *AbookClubScrapper) GetType() string {
	return "abook-club"
}

func convertNodeToBook(el *goquery.Selection) (*abooks.ABook, error) {
	a := el.Find("div.entry_header_full a")

	link, _ := a.Attr("href")
	linkParts := strings.Split(link, "=")
	if len(linkParts) < 2 {
		return nil, fmt.Errorf("Can't parse book id from '%s'", link)
	}
	id := linkParts[1]

	rawTitle := a.Text()
	rawTitleParts := strings.SplitN(rawTitle, "-", 2)
	var title, author string
	if len(rawTitleParts) == 2 {
		author = strings.Trim(rawTitleParts[0], " ")
		title = strings.Trim(rawTitleParts[1], " ")
	}

	authorAndTime := el.Find("div.entry_time").Text()
	authorAndTimeParts := strings.Split(authorAndTime, ", ")
	entryContent := el.Find("div.entry_content")
	description, _ := entryContent.Html()

	props := make(map[string]string)
	terms := strings.Split(strings.ReplaceAll(description, "<br/>", "\n"), "<b>")
	for _, t := range terms {
		termParts := strings.Split(t, "</b>")
		if len(termParts) == 2 {
			key := strings.Trim(termParts[0], " \n:")
			value := strings.Trim(termParts[1], " \n:")

			props[key] = value
		}
	}

	artists := []string{}

	if val, e := props["Исполнитель"]; e {
		for _, artist := range strings.Split(val, ";") {
			artists = append(artists, strings.Trim(artist, " "))
		}
	}

	var length int

	if val, e := props["Длительность"]; e {
		length := 0
		parts := strings.Split(val, ":")
		for _, part := range parts {
			partValue, err := strconv.Atoi(part)
			if err != nil {
				length = 0
				break
			}
			length = length*60 + partValue
		}
	}

	var size string
	if val, e := props["Размер"]; e {
		size = val
	}

	var quality string
	if val, e := props["Качество"]; e {
		quality = val
	}

	var desc string
	if val, e := props["Описание"]; e {
		desc = val
	}

	var year int
	if val, e := props["Год издания"]; e {
		year, _ = strconv.Atoi(val)
	}

	datetime := strings.Split(authorAndTimeParts[1], "\n")[0]

	created, err := time.Parse(timeFormat, datetime)

	if err != nil {
		return nil, fmt.Errorf("Can't parse book date %w", err)
	}

	book := &abooks.ABook{
		Id:          fmt.Sprintf("abook-club-%s", id),
		RawTitle:    rawTitle,
		Title:       title,
		Author:      author,
		Date:        created,
		Link:        link,
		Description: desc,
		Artists:     artists,
		Length:      length,
		Size:        size,
		Quality:     quality,
		Props:       props,
		Year:        year,
	}

	return book, nil
}

func parseBody(reader io.Reader) ([]abooks.ABook, error) {

	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	books := []abooks.ABook{}

	entries := document.Find("div.entry")

	for i := range entries.Nodes {
		el := entries.Eq(i)
		book, err := convertNodeToBook(el)
		if err != nil {
			return nil, err
		}
		books = append(books, *book)
	}

	if len(books) == 0 {
		return nil, fmt.Errorf("No books founded")
	}

	return books, nil
}
