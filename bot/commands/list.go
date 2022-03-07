package commands

import (
	"fmt"
	"bookrawl/app/dao"
	"strings"
)


type ListCommand struct {
	DaoHolder *dao.DaoHolder
}

func (cmd *ListCommand) GetName() string {
	return "list"
}

func (cmd *ListCommand) GetDescription() string {
	return "List of last processed books"
}

func (cmd *ListCommand) Run(ctx *Context) error{

	result, err := cmd.DaoHolder.GetBookStore().Find(nil, 20)
	if err != nil {
		return err
	}

	lines := []string{}

	for _, book := range result.Books {
		lines = append(lines, fmt.Sprintln(book.Author, "-", book.Title, "-", book.Date, book.AuthorId, book.Link))
	}

	ctx.Reply(strings.Join(lines, "\n"))

	return nil
}
