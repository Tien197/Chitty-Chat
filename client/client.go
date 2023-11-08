package main

import (
	"bufio"
	"context"
	"flag"
	"github.com/Tien197/Chitty-Chat/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"unicode/utf8"
)

type Client struct {
	proto.UnimplementedParticipantServiceServer // Necessary
	id                                          int
	portNumber                                  int
	lamportTime                                 int
}

var (
	clientPort = flag.Int("cPort", 0, "client port number")
	serverPort = flag.Int("sPort", 0, "server port number (should match the port used for the server)")
	clientID   = flag.Int("id", 0, "client ID number")
)

func main() {
	// Parse the flags to get the port for the client
	flag.Parse()

	// Create a client
	client := &Client{
		id:          *clientID,
		portNumber:  *clientPort,
		lamportTime: 1,
	}
	// Starts the client
	go startClient(client)

	// Wait for the client (user) to ask for the time
	go waitForJoinRequest(client)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT)

	// Block until a signal is received
	<-sigChan

	client.lamportTime++
	log.Printf("Client %d disconnected from the server at Lamport time %d", client.id, client.lamportTime)
}

func startClient(client *Client) {

	// Create a new grpc server
	grpcServer := grpc.NewServer()

	// Make the server listen at the given port (convert int port to string)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(client.portNumber))

	if err != nil {
		log.Fatalf("Could not create the server %v", err)
	}
	// Register the grpc server and serve its listener
	proto.RegisterParticipantServiceServer(grpcServer, client)

	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}
}

func waitForJoinRequest(client *Client) {
	// Connect to the server
	serverConnection, _ := connectToServer(client)

	client.lamportTime++
	log.Printf("Client %d requests to join server at Lamport Time %d", client.id, client.lamportTime)

	_, err := serverConnection.ParticipantJoins(context.Background(), &proto.ClientInfo{
		ClientId:    int64(client.id),
		LamportTime: int64(client.lamportTime),
		PortNumber:  int64(client.portNumber),
	})

	if err != nil {
		log.Printf("Client %d could not join server", client.id)
	}
	if err != nil {
		log.Printf("Error in ParticipantJoins: %v", err)
	}

	// Wait for input in the client terminal
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()

		if utf8.ValidString(input) && len(input) <= 128 {
			log.Printf("Client %d publishes message: \"%s\" at Lamport Time ...\n", client.id, input)
		} else {
			log.Print("Not a valid message! Send a message of UTF-8 and within 128 characters in length.")
		}

		// Ask the server for the time
		clientReturnMessage, err := serverConnection.ParticipantMessages(context.Background(), &proto.ClientInfo{
			ClientId: int64(client.id),
			Message:  input,
		})

		if err != nil {
			log.Printf(err.Error())
		} else {
			log.Printf("%s broadcasts client {id}'s \"%s\" at Lamport Time ...\n", clientReturnMessage.ServerName, input)
		}
	}
}

func connectToServer(client *Client) (proto.CCServiceClient, error) {
	// Dial the server at the specified port.
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(*serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to port %d", *serverPort)
	} else {
		log.Printf("Client %d connected to the server at port %d at Lamport Time %d\n", client.id, *serverPort, client.lamportTime)
	}
	return proto.NewCCServiceClient(conn), nil
}

// could be refactored
func (client *Client) ClientJoinReturn(ctx context.Context, in *proto.ClientInfo) (*proto.ServerInfo, error) {

	if client.lamportTime < int(in.LamportTime) {
		client.lamportTime = int(in.LamportTime)
	}
	client.lamportTime++

	log.Printf("Client %d joined at lamport timestamp %d\n", in.ClientId, client.lamportTime)

	return &proto.ServerInfo{
		LamportTime: int64(client.lamportTime),
	}, nil
}
