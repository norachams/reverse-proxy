## Introduction
This is a take-home test for Cohere, involving the development of an HTTP Reverse Proxy with specific functionality. It's a Go-based application designed to serve as an intermediary between client applications and origin servers. Included is the implementation of middleware to enforce a global in-flight request limit, which ensures that the proxy never exceeds the defined maximum number of concurrent requests to the origin server.
### What is a Reverse Proxy?
An HTTP reverse proxy is a server that sits between client devices and backend web servers. It receives incoming requests from clients and forwards them to the appropriate backend server. 
Reverse proxies are essential for load balancing, caching, security, and optimizing web application performance.

## Getting Started
### Prerequisites
Before using the HTTP Reverse Proxy, ensure that you have the following prerequisites installed:
- Go (Programming Language)
- A running Go development environment
- Knowledge of HTTP, proxies, and networking concepts [(or just use the resources I used to learn)](#resources-used)

### Building and Running
```bash
# build the server 
go build original-server/server.go
# run the server 
./server
# you should see the output 
'origin server listening on port 8081'

# in another terminal build the reverse-proxy
go build
# run the reverse-proxy
./reverse-proxy 
# you should see the output
'reverse proxy server listening on port 8080'

# in a separate terminal, test the proxy server and the original server using 
curl localhost:8080
#note: the 'received request' for the proxy is commented out on line 27
#I commented it out so it doesn't flood the testing

# to run the test run  
go test -count=1 -v ./...
```


## Resources Used
Prior to this project, I was new to HTTP reverse proxy development. I conducted extensive research to learn, build, and implement the necessary functionality. Below are the key resources I utilized:

#### Learning Resources:
1. [Cloudflare's Blog Post](https://www.cloudflare.com/learning/cdn/glossary/reverse-proxy/): A comprehensive blog post from Cloudflare that provides insights into the concept of reverse proxies.

2. [Traefik Labs' YouTube Video](https://www.youtube.com/watch?v=tWSmUsYLiE4&ab_channel=TraefikLabs): A YouTube video by Traefik Labs that visually explains the fundamentals of reverse proxies.

3. [Golang Redis Caching Guide](https://voskan.host/2023/08/14/golang-redis-caching/): A helpful guide on implementing Redis caching in Go, a valuable resource for optimizing caching in your proxy.

#### Code and Implementation Guides:
1. [Dev.to Tutorial](https://dev.to/b0r/implement-reverse-proxy-in-gogolang-2cp4): A detailed tutorial on implementing a reverse proxy in Go (Golang).

2. [Tpeczek's Blog Post](https://www.tpeczek.com/2017/08/implementing-concurrent-requests-limit.html): A blog post by Tpeczek that discusses implementing concurrent request limits, a key feature of the project.

3. [Go.dev Best Practices](https://go.dev/talks/2013/bestpractices.slide#1): Best practices for Go (Golang) development, which provided guidance on writing efficient and maintainable code.

4. [Go.dev Documentation](https://go.dev/doc/): Official Go (Golang) documentation for in-depth reference on the Go programming language.

#### Testing Resources:
1. [Golang.cafe Tutorial](https://golang.cafe/blog/golang-httptest-example.html): A tutorial on performing HTTP testing in Go using the httptest package.

2. [Cybertec's Blog Post](https://til.cybertec-postgresql.com/post/2019-11-07-How-to-turn-off-test-caching-for-golang/): A blog post that explains how to turn off test caching for Golang, which was useful for testing the project.


## Design Decisions 

**Functionality decision:** 
When presented with three possible functionality options to implement, I made the decision to pursue the "Global In-flight Request Limit" because it appealed to me and offered an engaging coding and learning experience.

**Global in-flight request limit implementation:**
The reverse proxy server ```proxy.go``` has been designed to control the number of concurrent requests made to the origin server. It limits the concurrent requests to a maximum of 10 (this can be adjusted). 

**Concurrency Control:**
To control concurrency, a synchronization lock ```sync.Mutex``` is used to prevent race conditions when updating the count of concurrent requests. When a request comes in, it checks if the current number of concurrent requests exceeds the limit.

**Request Forwarding:**
Forwards incoming client requests to the origin server specified by the ```originServerURL``` while preserving the request's Host, URL, and Request URI. The response from the origin server is then returned to the client.

**Simulating Slow Origin Server Response:**
In the ```server.go``` we introduce a delay in the origin server's response by using the ```time.Sleep``` function for 20 milliseconds. The delay allows us to test and verify the reverse proxy server's ability to enforce the specified concurrency limit (in this case, a maximum of 10 concurrent requests). It ensures that the server handles concurrent requests reaching the limit effectively. 

**Testing:**
The ```proxy_test.go``` file includes unit tests for the reverse proxy server. It tests the server's ability to handle both successful requests and requests that exceed the maximum allowed concurrency limit. It tracks the total requests sent to the server, reports the requests allowed within the concurrency limit, and distinguishes between successful and failed requests. It then provides a numerical count of both passed and failed requests. You can see the test output below.
```
=== RUN   TestContext_NewHandler
Total concurrent requests:  12
Total concurrent requests allowed:  10
request failed for concurrency limit
request failed for concurrency limit
request successful
request successful
request successful
request successful
request successful
request successful
request successful
request successful
request successful
request successful
Total Passed:  10
Total Failed:  2
--- PASS: TestContext_NewHandler (0.03s)
PASS
```
### How it works:
```server.go:```
- This is the main application file.
- It sets up a simple origin server on port 8081.
- The server sleeps for 20 milliseconds to simulate a slow response.
- It responds with "origin server response" for any incoming request.

```proxy.go:```
- This file contains the reverse proxy server logic.
- It limits the number of concurrent requests to the origin server using a concurrency limit of 10 (NoOfConcurrentRequests).
- When a request comes in, it checks if the current number of concurrent requests exceeds the limit. If it does, the server returns a "503 Service Unavailable" response.
- It uses a synchronization lock to ensure thread safety when updating the count of concurrent requests.
- The client's request is forwarded to the origin server while maintaining the request's Host, URL, and Request URI.
- The response from the origin server is returned to the client.

```proxy_test.go:```
- This file contains unit tests for the reverse proxy server.
- It tests the server's ability to handle both successful requests and requests that exceed the maximum allowed concurrency limit (NoOfConcurrentRequests).

### Limitations:
- **Fixed Request limit:** The code sets a fixed limit of 10 concurrent requests to the origin server using the NoOfConcurrentRequests constant. To change the limit change the value in the ```proxy.go``` file to the value you want. 
```go
const (
    // NoOfConcurrentRequests No more than 10 concurrent requests to the origin server
    NoOfConcurrentRequests = 10
)
```

- **Focused Testing:** The testing code primarily focuses on concurrency testing, ensuring that the proxy handles concurrent requests as expected. However, the test suite would benefit from expansion to cover more scenarios, including edge cases and error handling. 
    - **Edge cases:** The edge cases can encompass testing scenarios where the origin server is slow to respond or where it returns error codes. 

- **HTTP-Only Support:** The proxy implemented in the code only supports serving HTTP requests and responses, not HTTPS. While this is suitable for basic use cases, it might not meet security and encryption requirements in more complex scenarios. 

- **Lack of Request and Response Headers Handling:** The code currently does not demonstrate handling of request and response headers. Handling headers is required for authentication tokens, content negotiation, or routing decisions.


# Future work: Scaling

- **Sharded Rate Limiting:** As mentioned in the challenge description as one of the possible functionalities to implement we can limit the request rate. This entails tracking the number of requests from each IP address and, when it exceeds the limit, rejecting additional requests. The same concept applies to rate limiting based on header values. We can achieve this by increasing the count for each IP or header and decreasing it when the request is completed. 

- **Request Retries:** Another functionality mentioned in the challenge description as one of the options to implement is retries. We could implement a mechanism for request retries when the origin server responds with a status code greater than or equal to 500. A for loop can be used for retry attempts, with a predefined number of retries and time intervals. If a successful response is received, the loop terminates. Here's an example of retry logic:

```go
for retries := 0; retries < 5; retries++ {
    // Make the request
    response, err := httpClient.Do(request)
    if err != nil {
        // Handle the error
        time.Sleep(5 * time.Second) // Wait for 5 seconds before the next attempt
    } else if response.StatusCode >= 500 {
        // Retry
    } else {
        // Break the loop as the request was successful
        break
    }
}
```

- **Caching:** For improved scalability, we can consider implementing caching using Redis. Instead of locally incrementing and decrementing values, we can utilize Redis as a centralized caching solution. This ensures that multiple proxy servers can share and synchronize rate-limiting information. The code snippet below demonstrates a simplified implementation:

```go
type Context struct {
    currentConcurrentRequests int
    // To remove race conditions among concurrent requests
    lock sync.Mutex
    originalServer *url.URL
    redisClient *redis.Client // Initialize the Redis client
}
```

- **Load Balancing:** We can enhance the proxy's capabilities by introducing load balancing. This feature allows for even distribution of traffic across multiple origin servers, enhancing both resilience and performance.

- **Auto-Scaling:** We could implement auto-scaling functionality to enable the proxy to dynamically adjust its capacity based on traffic volume. This ensures optimal resource allocation as traffic patterns fluctuate.


## Future work: Security

- **HTTPS Support:** Implement HTTPS support to encrypt traffic between clients and the proxy, as well as between the proxy and the origin server.

- **Authentication and Authorization:** Add robust authentication and authorization mechanisms to control access to the proxy. This can include token-based authentication, API key validation, or integration with identity providers.

- **DDoS Protection:** Introduce DDoS protection mechanisms to prevent abuse and unauthorized access to your proxy.












