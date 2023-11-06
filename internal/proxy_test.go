package internal

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
)

func TestContext_NewHandler(t *testing.T) {
	originServerURL, err := url.Parse("http://127.0.0.1:8081")
	if err != nil {
		log.Fatal("invalid origin server URL")
	}

	ctx := &Context{
		currentConcurrentRequests: 0,
		originalServer:            originServerURL,
	}

	var extraConcurrentRequests = 2
	var totalCall = extraConcurrentRequests + NoOfConcurrentRequests

	fmt.Println("Total concurrent requests: ", totalCall)
	fmt.Println("Total concurrent requests allowed: ", NoOfConcurrentRequests)

	var wg sync.WaitGroup
	// wait until all concurrent requests are finished
	wg.Add(totalCall)

	// count the response with 200
	statusOkCount := 0
	// count the response with 503. Is the request fails for concurrency limit, it will return 503
	statusServiceUnavailableCount := 0

	doCall := func() {
		defer wg.Done()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		ctx.NewHandler(w, req)
		res := w.Result()

		if res.StatusCode == http.StatusOK {
			fmt.Println("request successful")
			statusOkCount++
		}
		if res.StatusCode == http.StatusServiceUnavailable {
			fmt.Println("request failed for concurrency limit")
			statusServiceUnavailableCount++
		}
	}

	for i := 0; i < totalCall; i++ {
		go doCall()
	}

	wg.Wait()

	fmt.Println("Total Passed: ", statusOkCount)
	fmt.Println("Total Failed: ", statusServiceUnavailableCount)

	if statusOkCount != NoOfConcurrentRequests {
		t.Errorf("expected 10 status ok, got %d", statusOkCount)
	}

	if statusServiceUnavailableCount != extraConcurrentRequests {
		t.Errorf("expected 2 status service unavailable, got %d", statusServiceUnavailableCount)
	}
}
