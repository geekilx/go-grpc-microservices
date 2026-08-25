package main

import (
	"context"
	"log"

	blogv1 "github.com/geekilx/grpc-course/proto/blog/v1"
	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BlogService struct {
	blogv1.UnimplementedBlogServiceServer
}

func CreateBlogService() *BlogService {
	return &BlogService{}
}

func (s *BlogService) CreateBlog(ctx context.Context, req *blogv1.Blog) (*blogv1.BlogId, error) {

	item := BlogItem{
		AuthorId: req.AuthorId,
		Title:    req.Title,
		Content:  req.Content,
	}

	res, err := collection.InsertOne(ctx, item)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	oid, ok := res.InsertedID.(bson.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, "cannot convert to OID")
	}

	return &blogv1.BlogId{
		Id: oid.Hex(),
	}, nil

}

func (s *BlogService) ReadBlog(ctx context.Context, req *blogv1.BlogId) (*blogv1.Blog, error) {

	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "specifing ID is required")
	}

	oid, err := bson.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot parse ID")
	}

	data := &BlogItem{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(ctx, filter)

	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot find blog with ID provided")
	}

	log.Printf("blog with ID %v found", req.Id)

	return documentToBlog(data), nil

}

func (s *BlogService) UpdateBlog(ctx context.Context, req *blogv1.Blog) (*emptypb.Empty, error) {

	oid, err := bson.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot parse ID")
	}

	data := &BlogItem{
		AuthorId: req.AuthorId,
		Title:    req.Title,
		Content:  req.Content,
	}
	res, err := collection.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": data})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not update: %v", err)
	}

	if res.MatchedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "cannot find blog with Id")
	}

	return &emptypb.Empty{}, nil

}

func (s *BlogService) ListBlogs(req *emptypb.Empty, stream blogv1.BlogService_ListBlogsServer) error {

	ctx := stream.Context()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return status.Errorf(codes.Internal, "unknown internal error: %v", err)
	}

	for cursor.Next(ctx) {
		blog := &BlogItem{}

		if err := cursor.Decode(blog); err != nil {
			return status.Errorf(codes.Internal, "failed to retrive blog from database")
		}

		if err := stream.Send(documentToBlog(blog)); err != nil {
			return status.Errorf(codes.Internal, "failed to send blog")
		}

	}
	return nil
}

func (s *BlogService) DeleteBlog(ctx context.Context, req *blogv1.BlogId) (*emptypb.Empty, error) {

	oid, err := bson.ObjectIDFromHex(req.Id)
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "couldn't convert to oid")
	}

	res, err := collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "couldn't delete the blog: %v", err)
	}

	if res.DeletedCount == 0 {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "no blog was deleted")
	}

	return &emptypb.Empty{}, nil

}
