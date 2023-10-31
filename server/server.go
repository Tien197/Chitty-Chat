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
	proto.UnimplementedTimeAskServer // Necessary
	name                             string
	port                             int
	lamportTime                      int
}

// Used to get the user-defined port for the server from the command line
var port = flag.Int("port", 0, "server port number")

func main() {
	// Get the port from the command line when the server is run
	flag.Parse()

	// Create a server struct
	server := &Server{
		name:        "Chitty-Chat",
		port:        *port,
		lamportTime: 1,
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
	proto.RegisterTimeAskServer(grpcServer, server)

	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}
}

func (s *Server) AskForTime(ctx context.Context, in *proto.AskForTimeMessage) (*proto.TimeMessage, error) {
	log.Printf("Participant %d sends message: \"%s\" at Lamport time ... \n", in.ClientId, in.Message)
	return &proto.TimeMessage{
		ServerName: s.name,
	}, nil
}

// when client joins server
func (s *Server) ParticipantJoinsServer(ctx context.Context, in *proto.AskForTimeMessage) (*proto.TimeMessage, error) {
	// updates lamport time depending on client
	if s.lamportTime < int(in.LamportTime) {
		s.lamportTime = int(in.LamportTime)
	}
	s.lamportTime++
	log.Printf("Participant %d joins %s at Lamport time %d\n", in.ClientId, s.name, s.lamportTime)

	// need to be broadcast to all existing participants
	s.lamportTime++
	log.Printf("%s broadcasts join message to Participant %d at Lamport time %d", s.name, in.ClientId, s.lamportTime)

	return &proto.TimeMessage{
		LamportTime: int64(s.lamportTime),
	}, nil
}
