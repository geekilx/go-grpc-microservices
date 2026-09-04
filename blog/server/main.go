package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	blogv1 "github.com/geekilx/go-grpc-microservices/proto/blog/v1"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var blogCollection *mongo.Collection
var userCollection *mongo.Collection

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, falling back to system environment variables.")
	}

	caCert, err := os.ReadFile("grpc-cets/ca.pem")
	if err != nil {
		log.Fatalf("failed to read CA certificate: %v", err)
	}
	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		log.Fatalf("failed to append CA certificate")
	}

	serverCert, err := tls.LoadX509KeyPair("grpc-cets/server.pem", "grpc-cets/server-key.pem")
	if err != nil {
		log.Fatalf("failed to load server certificate: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    caPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://root:root@localhost:27017"
	}

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("failed to connect to mongodb: %v", err)
	}

	blogCollection = client.Database("blogdb").Collection("blog")
	userCollection = client.Database("blogdb").Collection("user")

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = userCollection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("failed to create unique index: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening on port :50051")

	mongoRepo := NewMongoRepository(blogCollection)
	userRepo := NewUserRepository(userCollection)

	s := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)),
		grpc.ChainUnaryInterceptor(AuthInterceptor, LogginInterceptor))
	blogService := CreateBlogService(mongoRepo, userRepo)
	blogv1.RegisterBlogServiceServer(s, blogService)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down server")
	s.GracefulStop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Disconnect(ctx); err != nil {
		log.Printf("error disconnecting mongo: %v", err)
	}

}
