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

type RutrackerApiClient struct {
}

func (c *RutrackerApiClient) GetTopicsRegTimeSorted(forumId int) ([]TopicCreateInfo, error) {
	url := fmt.Sprintf("http://api.rutracker.org/v1/static/pvc/f/%d", forumId)

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
