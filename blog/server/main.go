package main

import (
	"log"
	"net"

	blogv1 "github.com/geekilx/grpc-course/proto/blog/v1"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
)

var collection *mongo.Collection

func main() {

	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://root:root@localhost:27017"))
	if err != nil {
		log.Fatalf("failed to connect to mongodb: %v", err)
	}

	collection = client.Database("blogdb").Collection("blog")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening on port :50051")

	s := grpc.NewServer()
	blogService := CreateBlogService()
	blogv1.RegisterBlogServiceServer(s, blogService)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
