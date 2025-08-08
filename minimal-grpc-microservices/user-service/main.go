package main

import (
	"context"
	"log"
	"net"

	pb "user-service/pb/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// In-memory storage for simplicity
var users = make(map[int64]*pb.CreateUserResponse)
var userID int64 = 1

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := &pb.CreateUserResponse{
		Id:    userID,
		Name:  req.Name,
		Email: req.Email,
	}

	users[userID] = user
	userID++

	log.Printf("Created user: %+v", user)
	return user, nil
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, exists := users[req.Id]
	if !exists {
		return nil, grpc.Errorf(codes.NotFound, "User not found")
	}

	response := &pb.GetUserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}

	log.Printf("Retrieved user: %+v", response)
	return response, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})

	log.Println("User service running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
