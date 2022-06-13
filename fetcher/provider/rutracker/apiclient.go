package rutracker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type TopicCreateInfo struct {
	TopicId int
	Created time.Time
}

type PVCResponseFormat struct {
	TopicId []string `json:"topic_id"`
}

type PVCResponse struct {
	Format PVCResponseFormat      `json:"format"`
	Result map[string]interface{} `json:"result"`
}

type PostTextResponse struct {
	Result map[string]interface{} `json:"result"`
}

type TorTopicData struct {
	InfoHash       string  `json:"info_hash"`
	ForumId        int     `json:"forum_id"`
	PosterId       int     `json:"poster_id"`
	Size           float64 `json:"size"`
	RegTime        int     `json:"reg_time"`
	TorStatus      int     `json:"tor_status"`
	Seeders        int     `json:"seeders"`
	TopicTitle     string  `json:"topic_title"`
	SeederLastSeen int     `json:"seeder_last_seen"`
	DlCount        int     `json:"dl_count"`
}

type TorTopicDataResult struct {
	Result map[string]TorTopicData `json:"result"`
}

type RutrackerApiClient struct {
	baseUrl string
}

func NewRutrackerApiClient(baseUrl string) *RutrackerApiClient {
	return &RutrackerApiClient{
		baseUrl: baseUrl,
	}
}

func (c *RutrackerApiClient) GetTopicsRegTimeSorted(forumId int) ([]TopicCreateInfo, error) {
	url := fmt.Sprintf("%s/v1/static/pvc/f/%d", c.baseUrl, forumId)

	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	content := PVCResponse{}
	err = json.Unmarshal(body, &content)
	if err != nil {
		return nil, err
	}

	regPropIdx := -1
	for i, prop := range content.Format.TopicId {
		if prop == "reg_time" {
			regPropIdx = i
			break
		}
	}

	if regPropIdx < 0 {
		return nil, fmt.Errorf("No reg time property")
	}

	infos := make([]TopicCreateInfo, 0, len(content.Result))
	for strId, r := range content.Result {
		topicId, err := strconv.Atoi(strId)
		if err != nil {
			return nil, err
		}

		data := r.([]interface{})
		if len(data) <= regPropIdx {
			continue
		}
		createdUnixTime := data[regPropIdx].(float64)
		infos = append(infos, TopicCreateInfo{
			TopicId: topicId,
			Created: time.Unix(int64(createdUnixTime), 0),
		})
	}

	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Created.After(infos[j].Created)
	})

	return infos, nil
}

func (c *RutrackerApiClient) GetPostText(topicId int) (string, error) {
	url := fmt.Sprintf("%s/v1/get_post_text?by=topic_id&val=%d", c.baseUrl, topicId)
	res, err := http.Get(url)

	if err != nil {
		return "", err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	content := PostTextResponse{}
	err = json.Unmarshal(body, &content)
	if err != nil {
		return "", err
	}

	strId := strconv.Itoa(topicId)

	if text, e := content.Result[strId]; e && text != nil {
		return text.(string), nil
	}

	return "", fmt.Errorf("No topic text returned")
}

func (c *RutrackerApiClient) GetTorTopicData(topicId int) (*TorTopicData, error) {
	url := fmt.Sprintf("%s/v1/get_tor_topic_data?by=topic_id&val=%d", c.baseUrl, topicId)
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	content := TorTopicDataResult{}
	err = json.Unmarshal(body, &content)
	if err != nil {
		return nil, err
	}

	strId := strconv.Itoa(topicId)

	if data, e := content.Result[strId]; e {
		return &data, nil
	}

	return nil, fmt.Errorf("No topic data present")
}
