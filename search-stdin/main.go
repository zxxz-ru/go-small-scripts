package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// usage echo -e 'first\nsecond\nthird\n' | go run main.go
func countGo(s string) int {
	resp, err := http.Get(s)
	if err != nil {
		log.Printf("%v", err)
		return -1
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("%v", err)
		return -2
	}
	str := string(body[:])
	idx := strings.Count(str, "Go")
	if idx > 0 {
		return idx
	} else {
		return -3
	}
}

func main() {
	var total int
	input := os.Stdin
	buf := bufio.NewReader(input)
	for {
		line, err := buf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		if len(line) > 0 {
			line = bytes.TrimSpace(line)
			url := string(line[:])
			res := countGo(url)
			total = total + res
			fmt.Printf("Go string #%d: %s\n", res, line)
		}
	}
	fmt.Printf("Total: %d: \n", total)
}
