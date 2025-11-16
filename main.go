package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func getLinesChannel(fileData io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer fileData.Close()
		defer close(out)

		result := ""

		for {
			// prepare a slike of 8 bytes
			data := make([]byte, 8)

			// read a bunch of data in that slice
			n, err := fileData.Read(data)

			if n == 0 {
				break
			}

			if err != nil {
				fmt.Printf("Unexpected error %s", err)
			}

			// Slice to actually number of bytes populated
			// although we want to populat 8 bytes it could be the case less than 8 were available
			data = data[:n]

			if i := bytes.IndexByte(data, '\n'); i != -1 {
				result += string(data[:i])

				if len(result) > 0 {
					out <- strings.ReplaceAll(result, "\n", "")

				}
				result = ""
				data = data[i+1:]
			}

			result += string(data)
		}

		if i := len(result); i > 0 {
			out <- result
		}
	}()

	return out
}

func main() {
	fileData, err := os.Open("./assets/byte-data.txt")

	if err != nil {
		log.Fatal("Encountered error: ", err)
	}

	lines := getLinesChannel(fileData)

	count := 1
	for line := range lines {
		fmt.Println("Line " + fmt.Sprint(count) + ": " + line)
		count++
	}

}
