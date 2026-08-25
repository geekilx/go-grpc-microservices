package main

import (
	"context"
	"io"
	"log"

	blogv1 "github.com/geekilx/grpc-course/proto/blog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := blogv1.NewBlogServiceClient(conn)

	log.Println("Calling CreateBlog")

	res, err := client.CreateBlog(context.Background(), &blogv1.Blog{AuthorId: "ilia", Title: "Hello World", Content: "This is my first blog"})
	if err != nil {
		log.Fatalf("could not create blog: %v", err)
	}

	log.Printf("Created blog with id %s", res.Id)

	log.Println("Calling ReadBlog")

	rbRes, err := client.ReadBlog(context.Background(), &blogv1.BlogId{Id: "6a8a89b4f4ec8560e2c84b77"})
	if err != nil {
		log.Fatalf("failed to get blog: %v", err)
	}

	log.Printf("retrived blog is: %v", rbRes)

	log.Println("Calling UpdateBlog")

	_, err = client.UpdateBlog(context.Background(), &blogv1.Blog{Id: res.Id, AuthorId: "updated author", Content: "updated content", Title: "title updated"})
	if err != nil {
		log.Fatalf("failed to update blog: %v", err)
	}

	log.Println("updated blog successfully")

	log.Println("calling ListBlogs")

	blogStream, err := client.ListBlogs(context.Background(), nil)
	for {
		res, err := blogStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to receive from stream: %v", err)
		}

		log.Printf("received blog is: %+v", res)

	}

	log.Println("calling DeleteBlog")

	_, err = client.DeleteBlog(context.Background(), &blogv1.BlogId{Id: "6a8a89b4f4ec8560e2c84b77"})
	if err != nil {
		log.Fatalf("failed to delete blog: %v", err)
	}
	log.Println("blog was deleted successfully")

}
