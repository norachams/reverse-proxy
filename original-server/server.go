package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	originServerHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Printf("[origin server] received request at: %s\n", time.Now())

		fmt.Println("origin server sleeping for 20 milliseconds")
		// simulate a slow origin server.
		// This will help to simulate the reverse proxy server's concurrency limit
		time.Sleep(time.Millisecond * 20)

		_, _ = fmt.Fprint(rw, "origin server response\n")
	})

	fmt.Println("origin server listening on port 8081")
	log.Fatal(http.ListenAndServe(":8081", originServerHandler))
}
