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
	"time"
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
	log.Printf("Started server at port: %d at Lamport time %d \n", server.port, server.lamportTime)

	// Register the grpc server and serve its listener
	proto.RegisterTimeAskServer(grpcServer, server)

	// message sent to server terminal when client joins
	go func() { // should probably be refactored to get Lamport Time from client
		for {
			_, err := listener.Accept()
			if err != nil {
				log.Printf("Server failed to accept connection: %v", err)
				continue
			}
			log.Printf("Participant {id} joined %s at Lamport Time ... ", server.name)
		}
	}()

	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}
}

func (c *Server) AskForTime(ctx context.Context, in *proto.AskForTimeMessage) (*proto.TimeMessage, error) {
	log.Printf("Participant %d sends message: at Lamport time ... \n", in.ClientId)
	return &proto.TimeMessage{
		Time:       time.Now().String(),
		ServerName: c.name,
	}, nil
}

// when client joins server
/*
func (c *Server) ClientJoinsServer(ctx context.Context, in *proto.AskForTimeMessage) (*proto.TimeMessage, error) {
	log.Printf("Client %d joins %s at Lamport time ... \n", in.ClientId, c.name)
	return &proto.TimeMessage{
		Time:       time.Now().String(),
		ServerName: c.name,
	}, nil
}*/
