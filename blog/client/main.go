package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"os"
	"time"

	blogv1 "github.com/geekilx/grpc-course/proto/blog/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
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
	})),
		grpc.WithUnaryInterceptor(UnaryTimeoutInterceptor(5*time.Second)),
		grpc.WithStreamInterceptor(StreamTimeoutInterceptor(10*time.Second)),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := blogv1.NewBlogServiceClient(conn)

	log.Println("calling Signup")
	signRes, err := client.Signup(context.Background(), &blogv1.SignupRequest{Username: "ilia", Password: "pass", Email: "ilia@gmail.com"})
	if err != nil {
		log.Fatalf("failed to signup: %v", err)
	}
	log.Printf("user %v successfully signed up", signRes.Id)

	log.Println("calling Login")

	loginRes, err := client.Login(context.Background(), &blogv1.LoginRequest{Username: "ilia", Password: "pass"})
	if err != nil {
		log.Fatalf("failed to login: %v", err)
	}

	md := metadata.Pairs("Authorization", "Bearer "+loginRes.Token)

	authCtx := metadata.NewOutgoingContext(context.Background(), md)

	log.Printf("login token is: %v", loginRes.Token)

	log.Println("Calling CreateBlog")

	res, err := client.CreateBlog(authCtx, &blogv1.Blog{AuthorId: "ilia", Title: "Hello World", Content: "This is my first blog"})
	if err != nil {
		log.Fatalf("could not create blog: %v", err)
	}

	log.Printf("Created blog with id %s", res.Id)

	log.Println("Calling ReadBlog")

	rbRes, err := client.ReadBlog(authCtx, &blogv1.BlogId{Id: res.Id})
	if err != nil {
		log.Fatalf("failed to get blog: %v", err)
	}

	log.Printf("retrived blog is: %v", rbRes)

	log.Println("Calling UpdateBlog")

	_, err = client.UpdateBlog(authCtx, &blogv1.Blog{Id: res.Id, AuthorId: "updated author", Content: "updated content", Title: "title updated"})
	if err != nil {
		log.Fatalf("failed to update blog: %v", err)
	}

	log.Println("updated blog successfully")

	log.Println("calling ListBlogs")

	blogStream, err := client.ListBlogs(authCtx, nil)
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

	_, err = client.DeleteBlog(authCtx, &blogv1.BlogId{Id: res.Id})
	if err != nil {
		log.Fatalf("failed to delete blog: %v", err)
	}
	log.Println("blog was deleted successfully")

	log.Println("calling DeleteUser")

	_, err = client.DeleteUser(authCtx, &blogv1.DeleteUserRequest{Id: signRes.Id})
	if err != nil {
		log.Fatalf("failed to delete user: %v", err)
	}
	log.Println("user was deleted successfully")

}
