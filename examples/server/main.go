package main

import (
	"context"
	"log"
	"net"

	userpb "github.com/m0rk0vka/grpc/gen/go/user/v1"
	"google.golang.org/grpc"
)

type userService struct {
	userpb.UnimplementedUserServiceServer
}

func (u *userService) GetUser(_ context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	return &userpb.GetUserResponse{
		User: &userpb.User{
			Uuid:     req.Uuid,
			FullName: "Vladimir",
		},
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:9879")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	userpb.RegisterUserServiceServer(grpcServer, &userService{})

	grpcServer.Serve(lis)
}
