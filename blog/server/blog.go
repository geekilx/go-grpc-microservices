package main

import (
	"context"
	"errors"
	"os"
	"time"

	blogv1 "github.com/geekilx/go-grpc-microservices/proto/blog/v1"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BlogService struct {
	blogv1.UnimplementedBlogServiceServer
	blogRepo BlogRepository
	userRepo UserRepository
}

func CreateBlogService(blogRepo BlogRepository, userRepo UserRepository) *BlogService {
	return &BlogService{
		blogRepo: blogRepo,
		userRepo: userRepo,
	}

}

func (s *BlogService) CreateBlog(ctx context.Context, req *blogv1.Blog) (*blogv1.BlogId, error) {
	item := &BlogItem{
		AuthorId: GetUserID(ctx),
		Title:    req.Title,
		Content:  req.Content,
	}

	oid, err := s.blogRepo.Create(ctx, item)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	return &blogv1.BlogId{Id: oid.Hex()}, nil

}

func (s *BlogService) ReadBlog(ctx context.Context, req *blogv1.BlogId) (*blogv1.Blog, error) {

	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "specifing ID is required")
	}

	oid, err := bson.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot parse ID")
	}

	item, err := s.blogRepo.GetByID(ctx, oid)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "cannot find blog with ID provided")
		}

		return nil, status.Errorf(codes.Internal, "failed to retrive blog: %v", err)
	}

	return documentToBlog(item), nil

}

func (s *BlogService) UpdateBlog(ctx context.Context, req *blogv1.Blog) (*emptypb.Empty, error) {

	oid, err := bson.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot parse ID")
	}

	data := &BlogItem{
		AuthorId: GetUserID(ctx),
		Title:    req.Title,
		Content:  req.Content,
	}

	if err := s.blogRepo.Update(ctx, oid, GetUserID(ctx), data); err != nil {
		if errors.Is(err, ErrNotFound) {
			return &emptypb.Empty{}, status.Errorf(codes.NotFound, "cannot find blog with Id")
		}
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "failed to update blog: %v", err)
	}

	return &emptypb.Empty{}, nil

}

func (s *BlogService) ListBlogs(req *emptypb.Empty, stream blogv1.BlogService_ListBlogsServer) error {

	ctx := stream.Context()

	items, err := s.blogRepo.List(ctx)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to list blogs: %v", err)
	}

	for _, item := range items {
		if err := stream.Send(documentToBlog(item)); err != nil {
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

	if err := s.blogRepo.Delete(ctx, oid, GetUserID(ctx)); err != nil {
		if errors.Is(err, ErrNotFound) {
			return &emptypb.Empty{}, status.Errorf(codes.NotFound, "no blog was deleted")
		}
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "failed to delete blog: %v", err)
	}

	return &emptypb.Empty{}, nil

}

func (s *BlogService) Signup(ctx context.Context, req *blogv1.SignupRequest) (*blogv1.SignupResponse, error) {

	if req.Username == "" || req.Password == "" || req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "username, password and email are required")
	}

	hashPass, err := hashPassword(req.Password)
	user := &UserItem{
		Username: req.Username,
		Password: hashPass,
		Email:    req.Email,
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	oid, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &blogv1.SignupResponse{Id: oid.Hex()}, nil

}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func passwordMatches(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func (s *BlogService) Login(ctx context.Context, req *blogv1.LoginRequest) (*blogv1.LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "username and password are required")
	}

	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, status.Errorf(codes.Unauthenticated, "invalid username or password")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	if !passwordMatches(req.Password, user.Password) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid username or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, status.Errorf(codes.Internal, "JWT_SECRET is not set")
	}
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to sign token: %v", err)
	}

	return &blogv1.LoginResponse{
		Token: tokenString,
	}, nil

}
func (s *BlogService) DeleteUser(ctx context.Context, req *blogv1.DeleteUserRequest) (*emptypb.Empty, error) {

	userID, err := bson.ObjectIDFromHex(GetUserID(ctx))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot convert to OID")
	}

	if err := s.userRepo.Delete(ctx, userID); err != nil {
		if errors.Is(err, ErrNotFound) {
			return &emptypb.Empty{}, status.Errorf(codes.NotFound, "user not found")
		}
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &emptypb.Empty{}, nil
}
