package main

import (
	"bookrawl/app/dao"
	"bookrawl/app/dao/abooks"
	"bookrawl/fetcher/provider/rutracker"
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
	err := processForum(2388, daoHolder)
	if err != nil {
		return err
	}

	err = processForum(2387, daoHolder)
	if err != nil {
		return err
	}

	return nil
}

func processForum(fid int, daoHolder *dao.DaoHolder) error {
	scrapper := rutracker.NewRutrackerHtmlScrapper()
	from := 0
	pageSize := 100
	for true {
		log.Printf("Process from %d\n", from)
		books, err := scrapper.GetBookPage(fid, from, pageSize)
		if err != nil {
			return err
		}

		if len(books) == 0 {
			break
		}

		err = processBooks(books, daoHolder)
		if err != nil {
			return err
		}

		from += len(books)
	}

	return nil
}

func processBooks(books []abooks.ABook, daoHolder *dao.DaoHolder) error {
	for _, book := range books {
		log.Println("Process book", book.RawTitle)
		author, err := daoHolder.GetAuthorsStore().FindByName(book.Author)
		if err != nil {
			log.Println("Error fetch author for book", book.RawTitle, err)
		} else if author == nil {
			log.Printf("No author [%s] founded for book: %s", book.Author, book.RawTitle)
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
