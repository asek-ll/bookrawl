package main

import (
	"bookrawl/app/abooks"
	"fmt"
	"sort"

	"context"
	"flag"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	connUri := flag.String("mongodbURI", "", "mongodb connection uri")
	flag.Parse()

	if *connUri == "" {
		log.Fatal("mongodbURI required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(*connUri))
	err = client.Ping(ctx, readpref.Primary())

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
	result, err := bookStore.Find(nil, 20)
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
	result, err := bookStore.Find(nil, 20)
	if err != nil {
		log.Fatal(err)
	}

	for _, book := range result.Books {
		fmt.Println(book.Author, "-", book.Title, "-", book.Date)
		//fmt.Println(book.RawTitle)
	}
}
