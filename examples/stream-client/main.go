package main

import (
	"context"
	"fmt"
	"io"
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	wearablepb "github.com/m0rk0vka/grpc-microservice-example/gen/go/wearable/v1"
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

	// some CLI to illustrate bits
	if err := ui.Init(); err != nil {
		log.Fatalf("Couldn't init UI: %v", err)
	}
	defer ui.Close()

	lc := widgets.NewPlot()
	lc.Title = "Heartbeat Per Second"
	lc.SetRect(0, 0, 70, 20)
	lc.Data = make([][]float64, 1)
	lc.DataLabels = []string{"Second", "Value"}
	lc.AxesColor = ui.ColorWhite
	lc.LineColors[0] = ui.ColorGreen
	lc.Marker = widgets.MarkerDot

	data := make([]float64, 60)
	ui.Render(lc)

	// cover and render new data

	res, err := client.BeatsPerSecond(context.Background(), &wearablepb.BeatsPerSecondRequest{
		Uuid: "Vladimir",
	})
	if err != nil {
		log.Fatalf("Couldn't request: %v", err)
	}

	go func() {
		for {
			resp, err := res.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Fatalln("Recieving", err)
			}

			data[int(resp.GetSecond())] = float64(resp.GetValue())
			lc.Data[0] = data

			ui.Render(lc)
		}
	}()

	uiEvents := ui.PollEvents()

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		case <-res.Context().Done():
			fmt.Println("All done, possible error", res.Context().Err())
			break
		}
	}
}
