package tasks

import "bookrawl/app/fantlab"

func fetchauthors(api fantlab.Api) error {
	resp, err := api.GetAllAuthors()

	if err != nil {
		return err
	}

	for _, author := range resp.List {
	}

	return nil
}
