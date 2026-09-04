package main

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrNotFound = errors.New("Document not foound")
)

type MongoBlogRepo struct {
	collection *mongo.Collection
}

func NewMongoRepository(collection *mongo.Collection) *MongoBlogRepo {
	return &MongoBlogRepo{
		collection: collection,
	}
}

func (r *MongoBlogRepo) Create(ctx context.Context, item *BlogItem) (bson.ObjectID, error) {
	res, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		return bson.NilObjectID, err
	}

	oid, ok := res.InsertedID.(bson.ObjectID)
	if !ok || oid == bson.NilObjectID {
		return bson.NilObjectID, errors.New("cannot convert to OID")
	}

	return oid, nil

}

func (r *MongoBlogRepo) GetByID(ctx context.Context, id bson.ObjectID) (*BlogItem, error) {
	filter := bson.M{"_id": id}
	item := &BlogItem{}

	if err := r.collection.FindOne(ctx, filter).Decode(item); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return item, nil
}

func (r *MongoBlogRepo) Update(ctx context.Context, id bson.ObjectID, author_id string, item *BlogItem) error {
	filter := bson.M{"_id": id, "author_id": author_id}
	update := bson.M{"$set": item}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil

}

func (r *MongoBlogRepo) Delete(ctx context.Context, id bson.ObjectID, author_id string) error {
	filter := bson.M{"_id": id, "author_id": author_id}
	res, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *MongoBlogRepo) List(ctx context.Context) ([]*BlogItem, error) {

	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	var items []*BlogItem
	for cursor.Next(ctx) {
		item := &BlogItem{}
		if err := cursor.Decode(item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil

}

type MongoUserRepo struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) *MongoUserRepo {
	return &MongoUserRepo{
		collection: collection,
	}
}

func (r *MongoUserRepo) Create(ctx context.Context, item *UserItem) (bson.ObjectID, error) {
	res, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return bson.NilObjectID, ErrNotFound
		}
		return bson.NilObjectID, err
	}

	oid, ok := res.InsertedID.(bson.ObjectID)
	if !ok || oid == bson.NilObjectID {
		return bson.NilObjectID, errors.New("cannot convert to OID")
	}

	return oid, nil

}

func (r *MongoUserRepo) GetByUsername(ctx context.Context, username string) (*UserItem, error) {
	filter := bson.M{"username": username}
	item := &UserItem{}

	if err := r.collection.FindOne(ctx, filter).Decode(item); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return item, nil
}

func (r *MongoUserRepo) Delete(ctx context.Context, id bson.ObjectID) error {
	filter := bson.M{"_id": id}
	res, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return ErrNotFound
	}

	return nil
}
