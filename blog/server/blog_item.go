package main

import (
	blogv1 "github.com/geekilx/grpc-course/proto/blog/v1"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type BlogItem struct {
	ID       bson.ObjectID `bson:"_id,omitempty"`
	AuthorId string        `bson:"author_id"`
	Title    string        `bson:"title"`
	Content  string        `bson:"content"`
}

func documentToBlog(data *BlogItem) *blogv1.Blog {

	return &blogv1.Blog{
		Id:       data.ID.Hex(),
		AuthorId: data.AuthorId,
		Title:    data.Title,
		Content:  data.Content,
	}

}
