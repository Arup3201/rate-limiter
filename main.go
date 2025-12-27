package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Arup3201/ratelimiter/algorithms"
	"github.com/Arup3201/ratelimiter/handlers"
)

const (
	HOST = "127.0.0.1"
	PORT = "8080"
)

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})

}

func main() {
	rateLimiter := algorithms.TokenBucketRateLimiter(5, 3)

	mux := http.NewServeMux()
	mux.Handle("GET /users", rateLimiter(http.HandlerFunc(handlers.GetUsers)))

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", HOST, PORT),
		Handler: logger(mux),

		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	log.Printf("server is starting at: %s:%s\n", HOST, PORT)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("server start error: %s\n", err)
	}
}
