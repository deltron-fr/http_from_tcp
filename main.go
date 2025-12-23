package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Println(err)
		return
	}

	// for concatenation of the read 8 bytes per each line.
	var currentLine strings.Builder

	for {
		rawData := make([]byte, 8)
		_, err = f.Read(rawData)
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Println(err)
			return
		}

		parts := strings.Split(string(rawData), "\n")
		for i := 0; i < len(parts); i++ {
			if len(parts) == 1 {
				currentLine.WriteString(parts[i])
			}
		}

		if len(parts) == 1 {
			continue
		}

		currentLine.WriteString(parts[0])
		fullLine := currentLine.String()
		fmt.Printf("read: %s\n", fullLine)

		currentLine.Reset()
		currentLine.WriteString(parts[1])
	}
}