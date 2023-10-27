package main

import (
	"math/rand"
	"time"

	wearablepb "github.com/m0rk0vka/grpc-microservice-example/gen/go/wearable/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type wearableService struct {
	wearablepb.UnimplementedWearableServiceServer
}

func (w *wearableService) BeatsPerSecond(req *wearablepb.BeatsPerSecondRequest, stream wearablepb.WearableService_BeatsPerSecondServer) error {
	for {
		select {
		case <-stream.Context().Done():
			return status.Error(codes.Canceled, "Stream has ended")
		default:
			time.Sleep(1 * time.Second)
			value := 40 + rand.Int31n(100)
			if err := stream.SendMsg(&wearablepb.BeatsPerSecondResponse{
				Value:  uint32(value),
				Second: uint32(time.Now().Second()),
			}); err != nil {
				return status.Error(codes.Canceled, "Stream has ended")
			}
		}
	}
}
