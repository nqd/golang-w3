package main

import (
	"errors"
	"log"
	"net"

	symonsayspb "github.com/nqd/golang-w3/simonsayspb"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	log.Println("Simon says hello")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	symonsayspb.RegisterSimonSaysServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) Game(stream symonsayspb.SimonSays_GameServer) (err error) {
	// ctx := stream.Context()

	req, err := stream.Recv()
	if err != nil {
		log.Printf("Error recieving: %v", err)
		return err
	}
	log.Println("Game function was invoked with a streaming request")
	log.Println(req)

	player := req.GetJoin()
	if player == nil {
		log.Printf("Error player is nil on initial join request: %v", req)
		err = errors.New("Player was nil on initial join request")
		return
	}
	// send back join request
	joinRes := &symonsayspb.Response{
		Event: &symonsayspb.Response_Turn{
			Turn: symonsayspb.Response_BEGIN,
		},
	}
	err = stream.SendMsg(joinRes)

	req, err = stream.Recv()
	if err != nil {
		log.Printf("Error recieving: %v", err)
		return err
	}

	color := req.GetPress()
	log.Printf("Color: %v", color)
	// send back turn request
	colorRes := &symonsayspb.Response{
		Event: &symonsayspb.Response_Lightup{
			Lightup: color,
		},
	}
	err = stream.SendMsg(colorRes)

	return
}
