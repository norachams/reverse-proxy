package internal

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	//"time"
)

const (
	// NoOfConcurrentRequests No more than 10 concurrent requests to the origin server
	NoOfConcurrentRequests = 10
)

type Context struct {
	currentConcurrentRequests int
	// To remove race condition among concurrent requests
	lock           sync.Mutex
	originalServer *url.URL
}

func (c *Context) NewHandler(rw http.ResponseWriter, req *http.Request) {
	//commented out so it doesnt flood the testing when running the test
	//fmt.Printf("[reverse proxy server] received request at: %s\n", time.Now())
	if c.currentConcurrentRequests >= NoOfConcurrentRequests {
		rw.WriteHeader(http.StatusServiceUnavailable)
		_, _ = fmt.Fprint(rw, "You have reached maximum concurrent requests. Please try again later.\n")
		return
	}

	// only one request can increase the number at a time
	c.lock.Lock()
	c.currentConcurrentRequests++
	c.lock.Unlock()

	// When the request is finished, decrease the number of concurrent requests
	defer func() {
		// only one request can decrease the number at a time
		c.lock.Lock()
		c.currentConcurrentRequests--
		c.lock.Unlock()
	}()

	// set req Host, URL and Request URI to forward a request to the origin server
	req.Host = c.originalServer.Host
	req.URL.Host = c.originalServer.Host
	req.URL.Scheme = c.originalServer.Scheme
	req.RequestURI = ""

	// save the response from the origin server
	originServerResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(rw, err)
		return
	}

	// return response to the client
	rw.WriteHeader(http.StatusOK)
	io.Copy(rw, originServerResponse.Body)
}

func SetProxy() {

	originServerURL, err := url.Parse("http://127.0.0.1:8081")
	if err != nil {
		log.Fatal("invalid origin server URL")
	}

	ctx := &Context{
		// current number of concurrent requests
		currentConcurrentRequests: 0,
		// original server URL
		originalServer: originServerURL,
	}

	reverseProxy := http.HandlerFunc(ctx.NewHandler)
	fmt.Println("reverse proxy server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", reverseProxy))
}
