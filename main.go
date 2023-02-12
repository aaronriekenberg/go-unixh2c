package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/net/http2"
)

const (
	unixSocketPath = "./socket"
)

func runConnectionHandler(
	conn net.Conn,
	handler http.Handler,
	http2Server *http2.Server,
) {
	log.Printf("begin runConnectionHandler")

	defer conn.Close()

	http2Server.ServeConn(
		conn,
		&http2.ServeConnOpts{
			Handler: handler,
		})

	log.Printf("end runConnectionHandler")
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	log.Printf("begin main unixSocketPath = %q", unixSocketPath)

	os.Remove(unixSocketPath)

	netListener, err := net.Listen("unix", unixSocketPath)
	if err != nil {
		log.Fatalf("net.Listen error path = %v: %v", unixSocketPath, err)
	}

	http2Server := &http2.Server{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("got http handler call r = %+v", r)
		fmt.Fprint(w, "Hello world")
	})

	for {
		conn, err := netListener.Accept()
		if err != nil {
			log.Fatalf("Accept error: %v", err)
		}

		go runConnectionHandler(conn, handler, http2Server)
	}
}
