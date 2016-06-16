package main

import (
	"fmt"
	"encoding/json"
	"net"
	"bufio"
)

type PeerInfo struct {
	Peer string `json:"peer"`
	Latency int `json:"latency"`
}

type UserInfo struct {
	Type string `json:"type" `
	User string `json:"user" `
	Room string `json:"room" `
	Host string `json:"host" `
	Latency []PeerInfo `json:"latency"`
}

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8888"
	CONN_TYPE = "tcp"
)

func main() {
	// Listen for incoming connections.
	listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	queue := make(chan UserInfo, 10) // Buffered channel with capacity of 10
	
	if err != nil {
		fmt.Println("Error listening:", err.Error())
	}
	
	// Close the listener when the application closes.
	defer listener.Close()
	
	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		
		// Handle connections in a new goroutine.
		go handleRequest(conn, queue)
		go handleTasks(conn, queue)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn, queue chan<- UserInfo) {
	defer conn.Close()
	
	input := bufio.NewScanner(conn)
	var userInfo UserInfo
	
	for input.Scan() {
		text := input.Text()
		byte_text := []byte(text)
		err := json.Unmarshal(byte_text, &userInfo)
		if err != nil {
			continue
		}
		queue <- userInfo // send userInfo to task queue
	}
}

func handleTasks(conn net.Conn, queue <-chan UserInfo) {
	for {
		var userInfo UserInfo
		userInfo = <- queue
		switch userInfo.Type {
			case "newUser": fmt.Println("newUser") 
			case "host": fmt.Println("host")
			case "disconnectedUser": fmt.Println("disconnectedUser")
		}
		fmt.Fprintf(conn, "Type: %s  User: %s  Room: %s  Host: %s", userInfo.Type, userInfo.User, userInfo.Room, userInfo.Host)
	}
}