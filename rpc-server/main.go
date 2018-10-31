package main

import (
	"com.github/zxxz/go/go-small-scripts/rpc-server/user"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"io"
	"net"
	"net/http"
	"net/rpc"
)

type HttpConn struct {
	in  io.Reader
	out io.Writer
}

func (c HttpConn) Read(p []byte) (n int, err error) {
	return c.in.Read(p)
}

func (c HttpConn) Write(d []byte) (n int, err error) {
	return c.out.Write(d)
}

func (c HttpConn) Close() error {
	return nil
}

func main() {
	server := rpc.NewServer()
	user := user.User{}
	server.Register(&user)

	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		panic(err)
	}

	defer listener.Close()

	http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/rpc" {
			serverCodec := jsonrpc2.NewServerCodec(&HttpConn{in: r.Body, out: w}, server)

			w.Header().Set("Content-type", "application/json")

			if err = server.ServeRequest(serverCodec); err != nil {
				http.Error(w, "Error while serving JSON request", 500)
				return
			}
		} else {
			http.Error(w, "Unknown request", 404)
		}
	}))
}
