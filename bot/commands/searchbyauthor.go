package commands

import (
	"bookrawl/app/dao"
	"bookrawl/app/dao/abooks"

	"fmt"
	"strings"
)

type SearchByAuthorCommand struct {
	DaoHolder *dao.DaoHolder
}

func (cmd *SearchByAuthorCommand) GetName() string {
	return "author"
}

func (cmd *SearchByAuthorCommand) GetDescription() string {
	return "List of last processed books"
}

func (cmd *SearchByAuthorCommand) Run(ctx *Context) error{
	authorName := ctx.Message.CommandArguments()
	astore := cmd.DaoHolder.GetAuthorsStore()
	author, err := astore.FindByName(authorName)
	if err != nil {
		return err
	}

	if author == nil {
		ctx.Reply(fmt.Sprintln("No author with name", authorName))
		return nil
	}


	bstore := cmd.DaoHolder.GetBookStore()

	filter := abooks.NewBooksFilterBuilder().SetAuthorId(&author.Id).Build()

	page, err := bstore.Find(filter, 20)
	if err != nil {
		return err
	}

	lines := []string{}

	for _, book := range page.Books {
		lines = append(lines, fmt.Sprintln(book.Author, "-", book.Title, "-", book.Date, book.AuthorId, book.Link))
	}

	if len(lines) == 0 {
		ctx.Reply(fmt.Sprintln("No books for author", author))
		return nil
	}

	ctx.Reply(strings.Join(lines, "\n"))

	return nil
}
