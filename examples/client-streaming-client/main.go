package main

import (
	"context"
	"fmt"
	"log"
	"time"

	wearablepb "github.com/m0rk0vka/grpc-microservice-example/gen/go/wearable/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	client := wearablepb.NewWearableServiceClient(conn)
	stream, err := client.ConsumerBeatPerSecond(context.Background())
	if err != nil {
		log.Fatalln("Consuming stream", err)
	}

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		if err := stream.Send(&wearablepb.ConsumerBeatPerSecondRequest{
			Uuid:   "Vladimir",
			Value:  uint32(i),
			Second: uint32(time.Now().Second()),
		}); err != nil {
			log.Fatalln("Sending value", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln("Closing", err)
	}

	fmt.Println("Total messages:", res.GetTotal())
}
