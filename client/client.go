package main

import (
	"bufio"
	"context"
	"flag"
	"github.com/Tien197/Chitty-Chat/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"unicode/utf8"
)

type Client struct {
	id          int
	portNumber  int
	lamportTime int
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

	// Wait for the client (user) to ask for the time
	go waitForTimeRequest(client)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT)

	// Block until a signal is received
	<-sigChan

	log.Printf("Client %d disconnected from the server", client.id)

}

func waitForTimeRequest(client *Client) {
	// Connect to the server
	serverConnection, _ := connectToServer(client)

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
		timeReturnMessage, err := serverConnection.AskForTime(context.Background(), &proto.AskForTimeMessage{
			ClientId: int64(client.id),
		})

		if err != nil {
			log.Printf(err.Error())
		} else {
			log.Printf("%s broadcasts client {id}'s {message} at Lamport Time ... %s\n", timeReturnMessage.ServerName, timeReturnMessage.Time)
		}
	}
}

func connectToServer(client *Client) (proto.TimeAskClient, error) {
	// Dial the server at the specified port.
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(*serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to port %d", *serverPort)
	} else {
		log.Printf("Client %d connected to the server at port %d at Lamport Time %d \n", client.id, *serverPort, client.lamportTime)
	}
	return proto.NewTimeAskClient(conn), nil
}

/*func (c *Client) ClientRequestToJoin(ctx context.Context, in *proto.TimeMessage) (*proto.AskForTimeMessage, error) {
	log.Printf("Client %d requested to join %s at Lamport time %d \n", c.id, in.ServerName, c.lamportTime+1)
	return &proto.AskForTimeMessage{
		ClientId: int64(c.id),
	}, nil
}*/
