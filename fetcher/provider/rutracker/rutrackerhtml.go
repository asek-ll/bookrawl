package rutracker

import (
	"bookrawl/app/dao/abooks"
	"bookrawl/fetcher/tasks"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var titleRegexp = regexp.MustCompile(`^(?P<author>.+)\s*[-–]\s*(?P<title>.+)\s*\[([^\]]*)\]\s*`)

type RutrackerHtmlScrapper struct {
	htmlClient HtmlClient
	apiClient  *RutrackerApiClient
}

func NewRutrackerHtmlScrapper() *RutrackerHtmlScrapper {
	return &RutrackerHtmlScrapper{
		htmlClient: HtmlClientImpl{},
		apiClient:  NewRutrackerApiClient("https://api.t-ru.org"),
	}
}

func (s *RutrackerHtmlScrapper) GetType() string {
	return "rutracker"
}

func (s *RutrackerHtmlScrapper) Fetch(task tasks.SyncTask, start time.Time) ([]abooks.ABook, error) {
	var forumId string
	exists := false
	if task.Params != nil {
		forumId, exists = task.Params["forumId"]
	}
	if !exists {
		return nil, fmt.Errorf("No forum id present")
	}

	fid, err := strconv.Atoi(forumId)
	if err != nil {
		return nil, err
	}

	return s.GetPageBooksCreatedAfter(fid, &task.LastRun, start)
}

func (s *RutrackerHtmlScrapper) GetBookPage(forumId int, from int, pageSize int) ([]abooks.ABook, error) {
	topics, err := s.apiClient.GetTopicsRegTimeSorted(forumId)

	if err != nil {
		return nil, err
	}

	books := make([]abooks.ABook, 0, pageSize)

	for i := from; i < len(topics) && i < from+pageSize; i += 1 {
		book, err := s.fetchBookFromTopic(topics[i])
		if err != nil {
			return nil, err
		}
		books = append(books, *book)
	}

	return books, nil
}

func (s *RutrackerHtmlScrapper) GetPageBooksCreatedAfter(forumId int, since *time.Time, start time.Time) ([]abooks.ABook, error) {
	topics, err := s.apiClient.GetTopicsRegTimeSorted(forumId)

	if err != nil {
		return nil, err
	}

	books := []abooks.ABook{}

	for _, topic := range topics {
		if since != nil && topic.Created.Before(*since) {
			break
		}
		if topic.Created.Before(start) {
			book, err := s.fetchBookFromTopic(topic)
			if err != nil {
				return nil, err
			}
			books = append(books, *book)
		}

	}

	return books, nil
}

func (s *RutrackerHtmlScrapper) fetchBookFromTopic(ti TopicCreateInfo) (*abooks.ABook, error) {
	topic, err := s.htmlClient.GetTopic(ti.TopicId)

	if err != nil {
		return nil, err
	}

	rawTitle := topic.Find("#topic-title").Text()

	var author string

	title := rawTitle

	titleParts := SplitAuthorTitleAndOther(title)
	if len(titleParts) >= 2 {
		author = titleParts[0]
		title = titleParts[1]
	}

	msg := topic.Find(".message").First()

	titleLink := msg.Find("p.post-time a")
	// posttime := titleLink.Text()
	link, _ := titleLink.Attr("href")
	body, _ := msg.Find("div.post_body").Html()

	propsHtml := strings.Split(body, "<span class=\"post-b\">")

	props := make(map[string]string)

	for i := 1; i < len(propsHtml); i += 1 {
		propHtml := propsHtml[i]
		kv := strings.Split(propHtml, "</span>:")
		if len(kv) != 2 {
			continue
		}
		key := kv[0]
		value := strings.TrimSuffix(strings.Trim(kv[1], "\n "), "<br/>")
		props[key] = value
	}

	var desc string
	if val, e := props["Описание"]; e {
		desc = val
	}

	artists := []string{}

	if val, e := props["Исполнитель"]; e {
		for _, artist := range strings.Split(val, ";") {
			artists = append(artists, strings.Trim(artist, " "))
		}
	}

	book := &abooks.ABook{
		Id:          fmt.Sprintf("rutracker-%d", ti.TopicId),
		RawTitle:    rawTitle,
		Title:       title,
		Author:      author,
		Link:        fmt.Sprintf("http://rutracker.org/forum/%s", link),
		Date:        ti.Created,
		Description: desc,
		Artists:     artists,
		// Length:      length,
		// Size:        size,
		// Quality:     quality,
		Props: props,
		// Year:        year,
	}

	return book, nil
}

func splitAuthorAndOther(title string) []string {

	seps := [][]string{
		{" - ", " – "},
		{"- ", "– "},
		{" -", " –"},
		{"-", "–"},
	}

	for _, ss := range seps {
		sepIndex := strings.Index(title, ss[0])
		anotherSepIndex := strings.Index(title, ss[1])
		if anotherSepIndex >= 0 && (anotherSepIndex < sepIndex || sepIndex < 0) {
			sepIndex = anotherSepIndex
		}

		if sepIndex >= 0 {
			author := strings.Trim(title[0:sepIndex], " ")
			title := strings.Trim(title[sepIndex:], " –-")
			return []string{author, title}
		}
	}

	return []string{title}
}

func SplitAuthorTitleAndOther(title string) []string {
	ps := splitAuthorAndOther(title)
	last := ps[len(ps)-1]

	propsStart := strings.Index(last, "[")

	if propsStart > 0 {
		title := strings.Trim(last[0:propsStart], " ")
		props := strings.Trim(last[propsStart+1:], " ")
		return append(ps[0:len(ps)-1], []string{title, props}...)
	}

	return ps
}
