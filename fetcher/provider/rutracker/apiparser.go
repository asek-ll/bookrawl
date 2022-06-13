package rutracker

import (
	"bookrawl/app/dao/abooks"
	"bookrawl/fetcher/tasks"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type RutrackerApiScrapper struct {
	apiClient *RutrackerApiClient
}

func NewRutrackerApiScrapper() *RutrackerApiScrapper {
	return &RutrackerApiScrapper{
		apiClient: NewRutrackerApiClient("https://api.t-ru.org"),
	}
}

func (s *RutrackerApiScrapper) GetType() string {
	return "rutracker"
}

func (s *RutrackerApiScrapper) Fetch(task tasks.SyncTask, start time.Time) ([]abooks.ABook, error) {
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

func (s *RutrackerApiScrapper) GetPageBooksCreatedAfter(forumId int, since *time.Time, start time.Time) ([]abooks.ABook, error) {
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

func (s *RutrackerApiScrapper) fetchBookFromTopic(ti TopicCreateInfo) (*abooks.ABook, error) {
	topic, err := s.apiClient.GetTorTopicData(ti.TopicId)

	if err != nil {
		return nil, err
	}

	body, err := s.apiClient.GetPostText(ti.TopicId)
	if err != nil {
		return nil, err
	}

	rawTitle := topic.TopicTitle

	var author string

	title := rawTitle
	fmt.Println(title)

	titleParts := SplitAuthorTitleAndOther(title)
	if len(titleParts) >= 2 {
		author = titleParts[0]
		title = titleParts[1]
	}

	propsHtml := strings.Split(body, "[b]")

	props := make(map[string]string)

	for i := 1; i < len(propsHtml); i += 1 {
		propHtml := propsHtml[i]
		kv := strings.Split(propHtml, "[/b]:")
		if len(kv) != 2 {
			continue
		}
		key := kv[0]
		value := strings.TrimSuffix(strings.Trim(kv[1], "\n "), "[hr]")
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
		Link:        fmt.Sprintf("https://rutracker.org/forum/viewtopic.php?t=%d", ti.TopicId),
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
