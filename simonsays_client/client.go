package main

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	symonsayspb "github.com/nqd/golang-w3/simonsayspb"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Simon says hello client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	defer cc.Close()

	c := symonsayspb.NewSimonSaysClient(cc)

	doBiDiStreaming(c)
}

func doBiDiStreaming(c symonsayspb.SimonSaysClient) {
	log.Println("Starting to do a BiDi Streaming RPC...")

	// we create a stream by invoking the client
	stream, err := c.Game(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	requests := []*symonsayspb.Request{
		{
			Event: &symonsayspb.Request_Join{
				Join: &symonsayspb.Request_Player{
					Id: "John",
				},
			},
		},
		// {
		// 	Event: &symonsayspb.Request_Join{
		// 		Join: &symonsayspb.Request_Player{
		// 			Id: "Jane",
		// 		},
		// 	},
		// },
		{
			Event: &symonsayspb.Request_Press{
				Press: symonsayspb.Color_BLUE,
			},
		},
		// {
		// 	Event: &symonsayspb.Request_Press{
		// 		Press: symonsayspb.Color_GREEN,
		// 	},
		// },
	}

	waitc := make(chan struct{})
	// we send a bunch of messages to the client (go routine)
	go func() {
		// function to send a bunch of messages
		for _, req := range requests {
			log.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// we receive a bunch of messages from the client (go routine)
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			log.Printf("Received: turn:%v, light:%v\n", res.GetTurn(), res.GetLightup())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}
