package main

import (
	"bookrawl/app/dao/abooks"
	"bookrawl/scheduler/utils"
	"fmt"
	"sort"

	"log"
)

func main() {
	client, err := utils.CreateMongoClient()

	if err != nil {
		log.Fatal(err)
	}
	bookStore := &abooks.AbookStore{
		Collection: client.Database("bookrawl").Collection("abooks"),
	}

	//printByAuthor(bookStore)
	printLast(bookStore)
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
