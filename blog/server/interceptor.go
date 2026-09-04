package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ctxKey string

const userIDKey ctxKey = "userID"

func GetUserID(ctx context.Context) string {
	id, _ := ctx.Value(userIDKey).(string)
	return id
}

func LogginInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("%s | %v | err=%v", info.FullMethod, time.Since(start), err)
	return resp, err
}

func AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

	switch info.FullMethod {
		case "/blog.v1.BlogService/Login", "/blog.v1.BlogService/Signup":
			return handler(ctx, req)
	}
	// if strings.HasSuffix(info.FullMethod, "/login") {
	// 	return handler(ctx, req)
	// }

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get("Authorization")) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing token")
	}

	tokenStr := strings.TrimPrefix(md.Get("Authorization")[0], "Bearer ")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "invalid claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "missing subject claim")
	}

	// Inject identity so handlers can access it
	ctx = context.WithValue(ctx, userIDKey, userID)

	return handler(ctx, req)
}
