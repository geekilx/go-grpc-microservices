package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"os"

	blogv1 "github.com/geekilx/grpc-course/proto/blog/v1"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var collection *mongo.Collection

func main() {

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

	s := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
	blogService := CreateBlogService()
	blogv1.RegisterBlogServiceServer(s, blogService)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
