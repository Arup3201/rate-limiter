package algorithms

import (
	"maps"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"time"
)

type tokenBucket struct {
	capacity   int // maximum tokens in the bucket - handles traffic bursts
	tokens     int // current number of tokens in the bucket
	rate       int // refill rate per second
	lastFilled time.Time
	retryAfter time.Duration
	mutex      sync.Mutex
}

func CreateTokenBucket(cap, rate int) *tokenBucket {
	return &tokenBucket{
		capacity: cap,
		tokens:   cap, // start with full bucket
		rate:     rate,
		mutex:    sync.Mutex{},
	}
}

func (tb *tokenBucket) refill() {
	now := time.Now()

	if tb.lastFilled.IsZero() {
		tb.lastFilled = now // initiate for next refill
		return
	}

	elapsed := now.Sub(tb.lastFilled)
	tokensToAdd := int(elapsed.Seconds()) * tb.rate

	if tokensToAdd > 0 {
		tb.lastFilled = now
		tb.tokens = min(tb.capacity, tb.tokens+tokensToAdd)
	}

	if tb.tokens == 0 && tokensToAdd == 0 {
		tb.retryAfter = time.Second - elapsed
	}
}

func (tb *tokenBucket) Take(tokens int) bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	tb.refill()

	if tokens <= tb.tokens {
		tb.tokens -= tokens
		return true
	}

	return false
}

type RateLimiter func(http.Handler) http.Handler

func TokenBucketRateLimiter(capacity, rate int) RateLimiter {
	bucket := CreateTokenBucket(capacity, rate)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if bucket.Take(1) {
				rec := httptest.NewRecorder()

				next.ServeHTTP(rec, r)

				maps.Copy(w.Header(), rec.Header())
				w.Header().Add("X-Ratelimit-Remaining", strconv.Itoa(bucket.tokens))
				w.Header().Add("X-Ratelimit-Limit", strconv.Itoa(bucket.capacity))
				w.WriteHeader(rec.Result().StatusCode)
				w.Write(rec.Body.Bytes())
			} else {
				w.Header().Add("X-Ratelimit-Limit", strconv.Itoa(bucket.capacity))
				w.Header().Add("X-Ratelimit-Retry-After", strconv.Itoa(int(bucket.retryAfter.Milliseconds())))
				w.WriteHeader(http.StatusTooManyRequests)
			}
		})
	}
}
