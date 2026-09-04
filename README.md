# Go gRPC Microservices & Blog Engine

A production-grade demonstration of building secure, resilient gRPC microservices in Go. This repository illustrates standard RPC communication patterns (unary, client streaming, server streaming, and bidirectional streaming), zero-trust Mutual TLS (mTLS), stateless JWT authentication via custom interceptors, database abstraction via the Repository Pattern, and containerization with Docker Compose.

---

## Architecture Overview

* **Transport & Protocol:** gRPC over HTTP/2 using Protocol Buffers (`proto3`).
* **Transport Security:** Mutual TLS (mTLS) with custom Root CA and client/server x509 certificates.
* **Authentication:** Stateless JSON Web Tokens (JWT) signed with HMAC-SHA256, validated through custom unary server interceptors.
* **Resilience:** Universal client-side context deadlines and timeout propagation via unary and stream interceptors.
* **Data Storage:** MongoDB (via the official `mongo-driver/v2`) decoupled from business logic using the Repository Pattern.
* **Lifecycle Management:** Graceful server shutdown listening for OS signals (`SIGINT`, `SIGTERM`) with clean connection termination.

---

## Services Implemented

1. **Blog Service (`proto/blog/v1/blog.proto`):**
   * **Authentication:** User registration (`Signup`) with `bcrypt` password hashing and session token generation (`Login`).
   * **CRUD Operations:** Authenticated `CreateBlog`, `ReadBlog`, `UpdateBlog`, `DeleteBlog`, and `DeleteUser`.
   * **Streaming:** Server streaming `ListBlogs` pipeline.
   * **Authorization:** Strict ownership validation on write/delete operations via context claims.

2. **Calculator Service (`proto/calculator/v1/calculator.proto`):**
   * Unary RPC (`Sum`)
   * Server Streaming (`Prime` factor decomposition)
   * Client Streaming (`Avg` running average calculation)
   * Bidirectional Streaming (`Max` rolling maximum tracking)

3. **Greet Service (`proto/greet/v1/greet.proto`):**
   * Baseline demonstrations of all four gRPC streaming paradigms.

---

## Project Structure

```text
├── blog/
│   ├── client/         # Authenticated blog gRPC client with interceptors
│   └── server/         # Blog service implementation, repositories, and auth
├── calculator/         # Calculator service (streaming paradigms)
├── docker/
│   ├── docker-compose.yml
│   └── dockerfile      # Multi-stage scratch build for the blog server
├── greet/              # Introductory gRPC greeting service
├── grpc-cets/          # CA, client, and server TLS certificates (development only)
├── proto/              # Protobuf definitions and generated Go code
├── go.mod
└── go.sum
