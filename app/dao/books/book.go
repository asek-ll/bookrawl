package books

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	Id        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Authors   []int              `bson:"authors"`
	FantLabId *int               `bson:"fantlabId"`
}
