package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	unixSocketPath = "./socket"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	log.Printf("begin main unixSocketPath = %q", unixSocketPath)

	os.Remove(unixSocketPath)

	netListener, err := net.Listen("unix", unixSocketPath)
	if err != nil {
		log.Fatalf("net.Listen error path = %v: %v", unixSocketPath, err)
	}

	h2s := &http2.Server{
		// ...
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("got http handler call r = %+v", r)
		fmt.Fprint(w, "Hello world")
	})

	h1s := &http.Server{
		Handler: h2c.NewHandler(handler, h2s),
	}

	log.Printf("before h1s.Serve")

	err = h1s.Serve(netListener)

	log.Fatalf("h1s.Serve error: %v", err)
}
