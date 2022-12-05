package fantlab

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Api struct {
}

type AuthorsResponse struct {
	List []Author `json:"list"`
}

type Author struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func NewApiClient() *Api {
	return &Api{}
}

func (api *Api) GetAllAuthors() (*AuthorsResponse, error) {
	response, err := http.Get("https://api.fantlab.ru/autorsall")

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var authors AuthorsResponse

	err = json.NewDecoder(response.Body).Decode(&authors)

	if err != nil {
		return nil, err
	}

	return &authors, nil
}

type Work struct {
	Id      int      `json:"work_id"`
	Title   string   `json:"title"`
	Name    string   `json:"work_name"`
	Authors []Author `json:"authors"`
}

func (api *Api) GetWork(workId int) (*Work, error) {
	response, err := http.Get(fmt.Sprintf("https://api.fantlab.ru/work/%d", workId))

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var work Work

	err = json.NewDecoder(response.Body).Decode(&work)

	if err != nil {
		return nil, err
	}

	return &work, nil
}
