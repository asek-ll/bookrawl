package fantlab

import (
	"encoding/json"
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
