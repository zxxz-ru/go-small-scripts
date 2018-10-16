package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

// usage echo -e 'first\nsecond\nthird\n' | go run main.go
func main() {
	input := os.Stdin
	buf := bufio.NewReader(input)
	idx := 0
	for {
		line, err := buf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				os.Exit(0)
			}
			log.Fatal(err)
		}

		line = bytes.TrimSpace(line)
		if len(line) > 0 {
			idx = idx + 1
			fmt.Printf("Line #%d: %s\n", idx, line)
		}
	}
}
