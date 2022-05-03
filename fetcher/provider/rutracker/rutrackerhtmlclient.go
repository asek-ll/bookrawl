package rutracker

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/charmap"
)

type HtmlClient interface {
	GetForum(forumid int, from int) (*goquery.Document, error)
	GetTopic(topicId int) (*goquery.Document, error)
}

type HtmlClientImpl struct {
}

func getDocument(url string) (*goquery.Document, error) {
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	decoder := charmap.Windows1251.NewDecoder()
	reader := decoder.Reader(response.Body)

	return goquery.NewDocumentFromReader(reader)
}

func (c HtmlClientImpl) GetForum(forumId int, from int) (*goquery.Document, error) {
	url := fmt.Sprintf("https://rutracker.org/forum/viewforum.php?f=%d&start=%d", forumId, from)

	return getDocument(url)
}

func (c HtmlClientImpl) GetTopic(topicId int) (*goquery.Document, error) {
	url := fmt.Sprintf("https://rutracker.org/forum/viewtopic.php?t=%d", topicId)
	return getDocument(url)
}
