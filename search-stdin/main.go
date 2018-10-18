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
func countGo(s string) chan int {
	ch := make(chan int, 1)
	go func() {
		defer close(ch)
		resp, err := http.Get(s)
		if err != nil {
			log.Printf("%v", err)
			ch <- -1
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("%v", err)
			ch <- -2
		}
		str := string(body[:])
		idx := strings.Count(str, "Go")
		if idx > 0 {
			ch <- idx
			fmt.Printf("Count for %s: %d\n", s, idx)
		} else {
			ch <- -3
		}
	}()
	return ch
}

func multyplexor(common, input chan int) {
	go func() {
		for i := range input {
			common <- i
		}
	}()
}

func main() {
	var total int
	var urls []string
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
			urls = append(urls, url)
		}
	}
	c := make(chan int, int(len(urls)/2))
	for _, url_string := range urls {
		multyplexor(c, countGo(url_string))
	}
	counter := 0
	for res := range c {
		counter = counter + 1
		if res > 0 {
			total = total + res
		}
		if counter == len(urls) {
			close(c)
		}
	}
	fmt.Printf("Total: %d: \n", total)
}
