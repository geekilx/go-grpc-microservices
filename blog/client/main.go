package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"os"

	blogv1 "github.com/geekilx/grpc-course/proto/blog/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	caCert, err := os.ReadFile("grpc-cets/ca.pem")
	if err != nil {
		log.Fatalf("failed to read CA certificate: %v", err)
	}
	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		log.Fatalf("failed to append CA certificate")
	}

	clientCert, err := tls.LoadX509KeyPair("grpc-cets/client.pem", "grpc-cets/client-key.pem")
	if err != nil {
		log.Fatalf("failed to load client certificate: %v", err)
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		ServerName:   "localhost",
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caPool,
	})))
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

	rbRes, err := client.ReadBlog(context.Background(), &blogv1.BlogId{Id: res.Id})
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
	if err != nil {
		log.Fatalf("failed while calling ListBlogs stream")
	}
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

	_, err = client.DeleteBlog(context.Background(), &blogv1.BlogId{Id: res.Id})
	if err != nil {
		log.Fatalf("failed to delete blog: %v", err)
	}
	log.Println("blog was deleted successfully")

}
