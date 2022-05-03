package main

import (
	"bookrawl/app/dao"
	"bookrawl/app/dao/abooks"
	"bookrawl/fetcher/provider/abookclub"
	"bookrawl/scheduler/utils"
	"log"
)

func main() {

	client, err := utils.CreateMongoClient()

	if err != nil {
		log.Fatal(err)
	}

	daoHolder := dao.NewDaoHolder(client)

	err = run(daoHolder)
	if err != nil {
		log.Fatal("Error on fetch all book", err)
	}
}

func run(daoHolder *dao.DaoHolder) error {
	abookScrapper := &abookclub.AbookClubScrapper{}

	pageNo := 1915
	for true {
		log.Println("Fetch page", pageNo)
		books, err := abookScrapper.GetPageBooks(pageNo)
		if err != nil {
			return err
		}

		if len(books) == 0 {
			break
		}

		err = processAbookPage(books, daoHolder)
		if err != nil {
			return err
		}

		pageNo += 1
	}

	return nil
}

func processAbookPage(books []abooks.ABook, daoHolder *dao.DaoHolder) error {
	for _, book := range books {
		log.Println("Process book", book.RawTitle)
		author, err := daoHolder.GetAuthorsStore().FindByName(book.Author)
		if err != nil {
			log.Println("Error fetch author for book", book, err)
		} else if author == nil {
			log.Println("No author founded for book", book)
		} else {
			book.AuthorId = []int{author.Id}
		}
	}

	err := daoHolder.GetBookStore().UpsertMany(books)

	if err != nil {
		return err
	}

	return nil
}
