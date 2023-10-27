package main

import (
	"context"
	"log"
	"net"
	"time"

	userpb "github.com/m0rk0vka/grpc-microservice-example/gen/go/user/v1"
	wearablepb "github.com/m0rk0vka/grpc-microservice-example/gen/go/wearable/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
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
	userServer := &userService{}
	wearableServer := &wearableService{}

	healthServer := health.NewServer()

	go func() {
		for {
			status := healthpb.HealthCheckResponse_SERVING
			// Check if user Service is valid
			if time.Now().Second()%2 == 0 {
				status = healthpb.HealthCheckResponse_NOT_SERVING
			}

			healthServer.SetServingStatus(userpb.UserService_ServiceDesc.ServiceName, status)
			healthServer.SetServingStatus("", status)

			time.Sleep(1 * time.Second)
		}
	}()

	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus(userpb.UserService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)

	userpb.RegisterUserServiceServer(grpcServer, userServer)
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	wearablepb.RegisterWearableServiceServer(grpcServer, wearableServer)

	grpcServer.Serve(lis)
}
