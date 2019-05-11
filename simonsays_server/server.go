package main

import (
	"fmt"
	"log"
	"net"

	symonsayspb "github.com/nqd/golang-w3/simonsayspb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct{}

func main() {
	fmt.Println("Simon says hello")
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

func (s *server) Game(stream symonsayspb.SimonSays_GameServer) error {
	return status.Errorf(codes.Unimplemented, "method Game not implemented")
}
