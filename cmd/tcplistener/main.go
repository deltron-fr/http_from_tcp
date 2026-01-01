package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}

		fmt.Println("Accepted connection from", conn.RemoteAddr())
		linesChan := getLinesChannel(conn)
		for line := range linesChan {
			fmt.Println(line)
		}
		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}

}

// getLinesChannel returns a channel that streams lines from the provided connection.
// It reads the messages in the connection in fixed-size byte chunks and reconstructs full lines that may span
// multiple reads by buffering partial data between iterations.
func getLinesChannel(c io.ReadCloser) <-chan string {
	// for concatenation of the read 8 bytes per each line.
	var currentLine strings.Builder

	lines := make(chan string)

	go func() {
		defer c.Close()
		defer close(lines)
		for {
			// reads data from connection chunk by chunk.
			rawData := make([]byte, 8)
			n, err := c.Read(rawData)
			if err != nil {
				if errors.Is(err, io.EOF) {
					if currentLine.Len() > 0 {
						lines <- currentLine.String()
					}
					break
				}
				log.Println(err)
				return
			}

			// splits 8 bytes data by the newline character
			// if there is no newline character, the length of the parts is 1
			// that means the end of the line has not been reached.
			parts := strings.Split(string(rawData[:n]), "\n")

			for i, part := range parts {
				if i < len(parts)-1 {
					currentLine.WriteString(part)
					lines <- currentLine.String()
					currentLine.Reset()
				} else {
					currentLine.WriteString(part)
				}
			}
		}
	}()

	return lines
}
