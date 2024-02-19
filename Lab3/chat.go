// Demonstration of channels with a chat application
// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Chat is a server that lets clients chat with each other.

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// Client struct holds the client's name and communication channel
type client struct {
	Name    string
	Channel chan<- string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[client]bool) // all connected clients

	for {
		select {
		case client := <-entering:
			clients[client] = true
			// Announce current clients to the new arrival
			go func() {
				client.Channel <- "Current clients:"
				for c := range clients {
					client.Channel <- c.Name
				}
			}()

		case client := <-leaving:
			delete(clients, client)

		case msg := <-messages:
			// Broadcast incoming message to all clients' outgoing message channels
			for client := range clients {
				client.Channel <- msg
			}
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	// Ask for client's name
	fmt.Fprint(conn, "Enter your name: ")
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	name := scanner.Text()

	// Create a new client with the connection and name
	Client := client{Name: name, Channel: ch}
	entering <- Client
	defer func() {
		leaving <- Client
		conn.Close()
	}()

	// Announce client's arrival
	messages <- name + " has arrived"

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- name + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()

	messages <- name + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}
