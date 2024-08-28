package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/go-sftp-server/gcs"
	"github.com/go-sftp-server/handler"
	"github.com/google/uuid"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var (
	User string
)

// Serve starts an SFTP server at the given address and port using the provided SSH configuration.
func Serve(sshConfig *ssh.ServerConfig, address, port string) {
	listenAt := fmt.Sprintf("%s:%s", address, port)
	listener, err := net.Listen("tcp", listenAt)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", listenAt, err)
	}
	defer listener.Close()

	log.Printf("Listening on %v\n", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn, sshConfig)
	}
}

// handleConnection handles incoming SSH connections in parallel.
func handleConnection(conn net.Conn, sshConfig *ssh.ServerConfig) {
	defer conn.Close()

	sconn, chans, reqs, err := ssh.NewServerConn(conn, sshConfig)
	if err != nil {
		log.Printf("failed to perform SSH handshake: %v", err)
		return
	}
	defer sconn.Close()

	User = fmt.Sprintf("%s: %s", sconn.User(), uuid.Must(uuid.NewRandom()).String())
	log.Printf("User: %s Logging in with bucket name: %s", sconn.User(), gcs.BucketName)

	// Discard incoming requests to keep the connection alive
	go ssh.DiscardRequests(reqs)

	// Handle incoming channels for the session
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Printf("User: %s Failed to accept channel: %v", User, err)
			continue
		}

		go handleRequests(requests)
		go handleSFTP(channel)
	}
}

// handleRequests handles SSH session requests, specifically the "subsystem" request for SFTP.
func handleRequests(requests <-chan *ssh.Request) {
	for req := range requests {
		if req.Type == "subsystem" && string(req.Payload[4:]) == "sftp" {
			req.Reply(true, nil)
		} else {
			req.Reply(false, nil)
		}
	}
}

// handleSFTP handles the SFTP subsystem for an accepted SSH channel.
func handleSFTP(channel ssh.Channel) {
	defer channel.Close()

	ctx := context.Background()
	sftpHandler, err := handler.Handler(ctx)
	if err != nil {
		log.Fatalf("User: %s failed to initialize gcs handler: %v", User, err)
	}

	server := sftp.NewRequestServer(channel, *sftpHandler)
	if err := server.Serve(); err == io.EOF {
		log.Printf("User: %s Sftp client exited session", User)
	} else if err != nil {
		log.Printf("User: %s Sftp server completed with error: %v", User, err)
	}
}
