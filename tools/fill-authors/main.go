package main

import (
	"bookrawl/app/dao"
	"bookrawl/app/dao/abooks"
	"bookrawl/fetcher/provider/rutracker"
	"bookrawl/scheduler/utils"
	"fmt"
	"log"
	"strings"
	"time"
)

func main() {

	client, err := utils.CreateMongoClient()

	if err != nil {
		log.Fatal(err)
	}

	daoHolder := dao.NewDaoHolder(client)

	err = run(daoHolder)
	// err = runSingle(daoHolder)
	if err != nil {
		log.Fatal("Error on fetch all book", err)
	}
}

func runSingle(daoHolder *dao.DaoHolder) error {
	abookStore := daoHolder.GetBookStore()
	book, err := abookStore.GetById("abook-club-90585")
	if err != nil {
		return err
	}

	fmt.Println(book.RawTitle)
	fmt.Println(rutracker.SplitAuthorTitleAndOther(book.RawTitle)[0])

	author, err := daoHolder.GetAuthorsStore().FindByName(book.Author)
	if err != nil {
		return err
	}

	fmt.Println(author)

	for _, s := range []string{
		"Марков-Бабкин Владимир - Новый Михаил 6, Император двух Империй [Юрий Белик, 2022, 60 kbps, MP3]",
	} {
		fmt.Println(strings.Join(rutracker.SplitAuthorTitleAndOther(s), "||"))
	}

	return nil
}

func run(daoHolder *dao.DaoHolder) error {
	abookStore := daoHolder.GetBookStore()
	authorStore := daoHolder.GetAuthorsStore()

	filter := abooks.NewBooksFilterBuilder().NoAuthor().Build()
	bookIdByAuthor := make(map[string][]string)
	cnt := 0

	for true {
		page, err := abookStore.Find(filter, 100)
		if err != nil {
			return err
		}
		var lastDate *time.Time = nil
		for _, book := range page.Books {

			if len(book.AuthorId) > 0 {
				continue
			}

			tp := rutracker.SplitAuthorTitleAndOther(book.RawTitle)
			if len(tp) >= 2 && tp[0] != book.Author {
				author, err := authorStore.FindByName(tp[0])
				if err != nil {
					return err
				}
				if author != nil {
					log.Printf("'%s' != '%s' from '%s'\n", tp[0], book.Author, book.RawTitle)
					log.Println(author, book)
					book.AuthorId = []int{author.Id}
					book.Author = tp[0]
					book.Title = tp[1]
					err = abookStore.Upsert(book)
					if err != nil {
						return err
					}
				}
			} else if strings.Index(book.Author, ",") > 0 {
				authors := strings.Split(book.Author, ",")
				resultAuthors := []int{}
				for _, a := range authors {
					single := strings.Trim(a, " ")
					author, err := authorStore.FindByName(single)
					if err != nil {
						return err
					}
					if author != nil {
						log.Println(author, book.Id)
						resultAuthors = append(resultAuthors, author.Id)
					}
				}
				if len(resultAuthors) > 0 {
					book.AuthorId = resultAuthors
					err = abookStore.Upsert(book)
					if err != nil {
						return err
					}
				}
			} else {
				// author, err := authorStore.FindByName(book.Author)
				// if err != nil {
				// 	return err
				// }
				// if author != nil {
				// 	log.Println("Find author", author, "for", book.RawTitle)
				// 	book.AuthorId = []int{author.Id}
				// 	err = abookStore.Upsert(book)
				// 	if err != nil {
				// 		return err
				// 	}
				// }

				v, e := bookIdByAuthor[book.Author]
				if !e {
					v = []string{book.Id}
				} else {
					v = append(v, book.Id)
				}
				bookIdByAuthor[book.Author] = v
			}
			lastDate = &book.Date
		}
		cnt += len(page.Books)
		log.Println("Process books", cnt)
		if lastDate == nil {
			break
		}
		filter = abooks.NewBooksFilterBuilder().SetBeforeDate(lastDate).NoAuthor().Build()
	}

	for k, v := range bookIdByAuthor {
		// log.Printf("auth: '%s', ids: %v\n", k, v)
		author, err := authorStore.FindByName(k)
		if err != nil {
			return err
		}
		if author != nil {
			log.Println("Founded", k, author, v)
			for _, id := range v {
				book, err := abookStore.GetById(id)
				if err != nil {
					return err
				}
				book.AuthorId = []int{author.Id}
				err = abookStore.Upsert(*book)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
