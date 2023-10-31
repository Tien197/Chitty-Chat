# Chitty-Chat

Distributed Systems (B-SWU & K-SD, Autumn 2023)
* Mandatory Hand-in 3 - Chitty Chat
* Submision Due Date: *`Tuesday 31, October 2023, 23:59`*

## How to run
Server has to run first in its own terminal, running in the root of the project.
```bash
go run server/server.go -port 5454
```

The clients are run in their own terminal as well. The id for the client can be changed.
```bash
go run client/client.go -cPort 8080 -sPort 5454 -id 1
```

## Authors
* Kasper Kirkegaard Nielsen (kkni@itu.dk)
* Omar Lukman Semou (omse@itu.dk)
* Tien Cam Ly (tily@itu.dk)

## Description:

You have to implement Chitty-Chat a distributed system, that is providing a chatting service, and keeps track of logical time using Lamport Timestamps.

We call clients of the Chitty-Chat service Participants. 


## System Requirements

1. Chitty-Chat is a distributed service, that enables its clients to chat. The service is using gRPC for communication. You have to design the API, including gRPC methods and data types. 

2. Clients in Chitty-Chat can Publish a valid chat message at any time they wish.  A valid message is a string of UTF-8 encoded text with a maximum length of 128 characters. A client **publishes** a message by making a gRPC call to Chitty-Chat.

3. The Chitty-Chat service has to broadcast every published message, together with the current logical timestamp, to all participants in the system, by using gRPC. It is an implementation decision left to the students, whether a Vector Clock or a Lamport timestamp is sent.

4. When a client receives a broadcasted message, it has to write the message and the current logical timestamp to the log

5. Chat clients can join at any time. 

6. A "Participant X  joined Chitty-Chat at Lamport time L" message is broadcast to all Participants when client X joins, including the new Participant.

7. Chat clients can drop out at any time. 

8. A "Participant X left Chitty-Chat at Lamport time L" message is broadcast to all remaining Participants when Participant X leaves.


## Technical Requirements:

1. Use gRPC for all messages passing between nodes

2. Use Golang to implement the service and clients

3. Every client has to be deployed as a separate process

4. Log all service calls (Publish, Broadcast, ...) using the [log](https://pkg.go.dev/log) package

5. Demonstrate that the system can be started with at least 3 client nodes 

6. Demonstrate that a client node can join the system

7. Demonstrate that a client node can leave the system

8. **Optional**: All elements of the Chitty-Chat service are deployed as Docker containers


## Hand-in requirements:

1. Hand in a single report in a pdf file

2. Discuss, whether you are going to use server-side streaming, client-side streaming, or bidirectional streaming? 

3. Describe your system architecture - do you have a server-client architecture, peer-to-peer, or something else?

4. Describe what  RPC methods are implemented, of what type, and what messages types are used for communication

5. Describe how you have implemented the calculation of the Lamport timestamps

6. Provide a diagram, that traces a sequence of RPC calls together with the Lamport timestamps, that corresponds to a chosen sequence of interactions: Client X joins, Client X Publishes, ..., 

7. Client X leaves. Include documentation (system logs) in your appendix.

8. Provide a link to a Git repo with your source code in the report
Include system logs, that document the requirements are met, in the appendix of your report


## Grading notes

* Partial implementations may be accepted, if the students can reason what they should have done in the report.

* The group must at least *attempt* to solve each of the requirements. The whole assignment is failed, if a requirement is not addressed.

