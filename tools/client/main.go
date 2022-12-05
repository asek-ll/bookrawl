package main

import (
	"bookrawl/app/dao"
	"bookrawl/app/dao/abooks"
	"bookrawl/app/dao/books"
	"bookrawl/app/fantlab"
	"bookrawl/scheduler/utils"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"

	"log"
)

func main() {
	client, err := utils.CreateMongoClient()

	if err != nil {
		log.Fatal(err)
	}

	daoHolder := dao.NewDaoHolder(client)
	command := os.Args[1]

	if command == "start" {
		workId, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}

		err = startBook(workId, os.Args[3], daoHolder)
	}

	if err != nil {
		log.Fatal(err)
	}

	//printByAuthor(bookStore)
	// printLast(bookStore)
}

func startBook(fantLabId int, abookId string, daoHolder *dao.DaoHolder) error {
	fantlabClient := fantlab.NewApiClient()
	work, err := fantlabClient.GetWork(fantLabId)
	if err != nil {
		return err
	}

	if work == nil {
		return errors.New("FantLab book not found with id")
	}

	bookStore := daoHolder.GetBookStore()
	book, err := bookStore.FindByFantLabId(fantLabId)

	if err != nil {
		return err
	}



	if book == nil {
		authorsIds := make([]int, len(work.Authors))
		for i, author := range work.Authors {
			authorsIds[i] = author.Id
		}
		book = &books.Book{
			Name:      work.Title,
			Authors:   authorsIds,
			FantLabId: &fantLabId,
		}
		err = bookStore.Create(book)
		if err != nil {
			return err
		}
	}

	
	stateStore := daoHolder.GetUserBookStateStore()

	stateStore.FindByBookIdAndUserId(book.Id, 


}

func printByAuthor(bookStore *abooks.AbookStore) {
	result, err := bookStore.Find(nil, 40)
	if err != nil {
		log.Fatal(err)
	}

	authors := []string{}
	titlesByAuthor := make(map[string][]string)

	for _, book := range result.Books {
		titles, e := titlesByAuthor[book.Author]
		if !e {
			titles = []string{book.Title}
			authors = append(authors, book.Author)
		} else {
			titles = append(titles, book.Title)
		}
		titlesByAuthor[book.Author] = titles
	}

	sort.Strings(authors)
	for i, author := range authors {
		fmt.Println(i, author)
		titles := titlesByAuthor[author]
		for _, title := range titles {
			fmt.Println("  *", title)
		}
	}

	bookCount := len(result.Books)
	if bookCount > 0 {
		fmt.Println("Last date", result.Books[bookCount-1].Date)
	}
}

func printLast(bookStore *abooks.AbookStore) {
	result, err := bookStore.Find(nil, 40)
	if err != nil {
		log.Fatal(err)
	}

	for _, book := range result.Books {
		fmt.Println(book.Author, "-", book.Title, "-", book.Date, book.AuthorId, book.Link)
		//fmt.Println(book.RawTitle)
	}
}
