package main

import (
	"fmt"
	"io"
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

func (w *wearableService) ConsumerBeatPerSecond(stream wearablepb.WearableService_ConsumerBeatPerSecondServer) error {
	var total uint32
	for {
		value, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&wearablepb.ConsumerBeatPerSecondResponse{
				Total: total,
			})
		}

		if err != nil {
			return err
		}

		fmt.Println(value.GetUuid(), value.GetValue(), value.GetSecond())
		total += 1
	}
}

func (w *wearableService) CalculatedBeatsPerSecond(stream wearablepb.WearableService_CalculatedBeatsPerSecondServer) error {
	var count, total uint32
	for {
		value, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		total += value.GetValue()
		count += 1
		if count%5 == 0 {
			avg := float32(total) / 5
			fmt.Println("Recieved", total, "sending", avg)
			if err := stream.Send(&wearablepb.CalculatedBeatsPerSecondResponse{
				Average: avg,
			}); err != nil {
				return err
			}
			total = 0
		}
	}
}
