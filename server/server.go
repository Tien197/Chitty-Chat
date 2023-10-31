package main

import (
	"context"
	"flag"
	"github.com/Tien197/Chitty-Chat/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

// Struct that will be used to represent the Server.
type Server struct {
	proto.UnimplementedClientToServerServer // Necessary
	name                                    string
	port                                    int
	lamportTime                             int
	participants                            []int
}

// Used to get the user-defined port for the server from the command line
var port = flag.Int("port", 0, "server port number")

func main() {
	// Get the port from the command line when the server is run
	flag.Parse()

	// Create a server struct
	server := &Server{
		name:         "Chitty-Chat",
		port:         *port,
		lamportTime:  1,
		participants: make([]int, 0),
	}

	// Start the server
	go startServer(server)

	// Keep the server running until it is manually quit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT)

	// Block until a signal is received
	<-sigChan

	log.Printf("%s was shut down", server.name)
}

func startServer(server *Server) {

	// Create a new grpc server
	grpcServer := grpc.NewServer()

	// Make the server listen at the given port (convert int port to string)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(server.port))

	if err != nil {
		log.Fatalf("Could not create the server %v", err)
	}
	log.Printf("Started %s at port: %d at Lamport time %d \n", server.name, server.port, server.lamportTime)

	// Register the grpc server and serve its listener
	proto.RegisterClientToServerServer(grpcServer, server)

	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}
}

// when participant sends message
func (s *Server) ParticipantMessages(ctx context.Context, in *proto.ClientInfo) (*proto.ServerInfo, error) {
	log.Printf("Participant %d sends message: \"%s\" at Lamport time ... \n", in.ClientId, in.Message)
	return &proto.ServerInfo{
		ServerName: s.name,
	}, nil
}

// when participant joins server
func (s *Server) ParticipantJoins(ctx context.Context, in *proto.ClientInfo) (*proto.ServerInfo, error) {
	// updates lamport time depending on participant

	if s.lamportTime < int(in.LamportTime) {
		s.lamportTime = int(in.LamportTime)
	}
	s.lamportTime++
	log.Printf("Participant %d joins %s at Lamport time %d\n", in.ClientId, s.name, s.lamportTime)

	// Assuming that all clientIds are unique
	s.participants = append(s.participants, int(in.ClientId))

	// need to be broadcast to all existing participants
	/*for participantId, participantPort := range s.participants {
		s.lamportTime++
		log.Printf("%s broadcasts join message to Participant %d at Lamport time %d", s.name, in.ClientId, s.lamportTime)
	}*/

	return &proto.ServerInfo{
		LamportTime: int64(s.lamportTime),
	}, nil
}
