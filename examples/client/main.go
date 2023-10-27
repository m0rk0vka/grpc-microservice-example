package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userpb "github.com/m0rk0vka/grpc-microservice-example/gen/go/user/v1"
)

func main() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial("localhost:9879", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	defer conn.Close()

	client := userpb.NewUserServiceClient(conn)

	res, err := client.GetUser(context.Background(), &userpb.GetUserRequest{
		Uuid: "hello",
	})
	if err != nil {
		log.Fatalf("fail to GetUser: %v", err)
	}

	fmt.Printf("%+v\n", res)
}
