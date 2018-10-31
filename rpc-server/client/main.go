package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	buf := bytes.NewReader([]byte(`{"jsonrpc":"2.0", "method":"User.Test", "params":{"Message":"It's alive!!"}, "id":"1"}`))
	res, err := http.Post("http://localhost:8080/rpc", "application/json", buf)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Fatal("ups")
	}
	io.Copy(os.Stdout, res.Body)
}
