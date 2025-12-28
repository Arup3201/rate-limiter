# Rate Limiter

Rate limiter in Golang.

## Token Bucket Algorithm

Example: A rate limiter using token bucket algorithm with bucket size=5 and refill rate per second=3.

```go
func main() {
    // ...
    rateLimiter := algorithms.TokenBucketRateLimiter(5, 3)

	mux := http.NewServeMux()
	mux.Handle("GET /users", rateLimiter(http.HandlerFunc(handlers.GetUsers)))

    // ...

}
```

Problem:

At the edge of timelines with 1 second interval there could be traffic burst and it will allow more than allowed number of requests to pass.

Example: Following log shows response code and the remaining number of requests.

```sh
2025/12/28 20:57:39.219252 200 Remaining: 4
2025/12/28 20:57:39.420753 200 Remaining: 3
2025/12/28 20:57:39.622352 200 Remaining: 2
2025/12/28 20:57:39.823871 200 Remaining: 1
2025/12/28 20:57:40.025374 200 Remaining: 0
2025/12/28 20:57:40.226914 200 Remaining: 2
2025/12/28 20:57:40.428263 200 Remaining: 1
2025/12/28 20:57:40.629600 200 Remaining: 0
2025/12/28 20:57:40.831237 429 Retry After: 395
2025/12/28 20:57:41.032134 429 Retry After: 194
```
