package main

import (
	// "bufio"
	// "bytes"
	"flag"
	"fmt"
	// "io"
	"io/ioutil"
	"log"
	"net/http"
	// "os"
	"regexp"
)

// usage echo -e 'first\nsecond\nthird\n' | go run main.go
func countGo(url string, re *regexp.Regexp) chan int {
	ch := make(chan int, 1)
	go func() {
		resp, err := http.Get(url)
		defer func() {
			close(ch)
			resp.Body.Close()
		}()
		if err != nil {
			log.Printf("%v", err)
			ch <- -1
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("%v", err)
			ch <- -2
			return
		}
		ma := re.FindAll(body, -1)
		if ma == nil {
			ch <- -3
			return
		}
		idx := len(ma)
		if idx > 0 {
			ch <- idx
			fmt.Printf("Count for %s: %d\n", url, idx)
			return
		} else {
			ch <- -4
			return
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
	// var urls []string
	var search_str string
	flag.StringVar(&search_str, "e", "Go", "regular expression to search in location")
	args := flag.Args()
	flag.Parse()
	fmt.Printf("%s\n", search_str)
	// input := ourl.Stdin
	// buf := bufio.NewReader(input)
	// for {
	// 	line, err := buf.ReadBytes('\n')
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		log.Fatal(err)
	// 	}
	//
	// 		if len(line) > 0 {
	// 			line = bytes.TrimSpace(line)
	// 			url := string(line[:])
	// 			urls = append(urls, url)
	// 		}
	// 	}
	if args == nil || len(args) == 0 {
		fmt.Println("Must provide at least one url for search")
		return
	}
	cb := int(len(args) / 2)
	if cb <= 0 {
		cb = 1
	}
	c := make(chan int, cb)
	defer close(c)
	re, err := regexp.Compile(search_str)
	if err != nil {
		fmt.Println("Can not compile provided regexp")
		return
	}
	for _, url_string := range args {
		multyplexor(c, countGo(url_string, re))
	}
	counter := 0
	for res := range c {
		counter = counter + 1
		if res > 0 {
			total = total + res
		}
		if counter == len(args) {
			break
		}
	}
	fmt.Printf("Total: %d: \n", total)
}
