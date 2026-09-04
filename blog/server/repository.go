package main

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BlogRepository interface {
	Create(ctx context.Context, item *BlogItem) (bson.ObjectID, error)
	GetByID(ctx context.Context, id bson.ObjectID) (*BlogItem, error)
	Update(ctx context.Context, id bson.ObjectID, author_id string, item *BlogItem) error
	Delete(ctx context.Context, id bson.ObjectID, author_id string) error
	List(ctx context.Context) ([]*BlogItem, error)
}

type UserRepository interface {
	Create(ctx context.Context, item *UserItem) (bson.ObjectID, error)
	GetByUsername(ctx context.Context, username string) (*UserItem, error)
	Delete(ctx context.Context, id bson.ObjectID) error
}
