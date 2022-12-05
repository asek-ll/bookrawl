package userbookstates

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Status int

const (
	Empty Status = iota
	WantToRead
	Reading
	Readed
)

type BookRead struct {
	AbookId primitive.ObjectID `bson:"abookId"`
	Start   time.Time          `bson:"start"`
	End     *time.Time         `bson:"end"`
}

type State struct {
	UserId primitive.ObjectID `bson:"userId"`
	BookId primitive.ObjectID `bson:"bookId"`
	Status Status             `bson:"status"`
	Rating int                `bson:"rating"`
	Reads  []BookRead         `bson:"reads"`
}
