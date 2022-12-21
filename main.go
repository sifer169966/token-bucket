package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	rateLimit = 10       // Requests per second
	bucketSize = rateLimit // Bucket size
)

type TokenBucket struct {
	capacity int           // Bucket capacity
	tokens   int           // Current number of tokens in the bucket
	rate     float64       // Rate at which tokens are added to the bucket (tokens/second)
	interval time.Duration // Interval at which tokens are added to the bucket
	mutex    sync.Mutex    // Mutex for synchronizing access to the bucket
	timer    *time.Timer   // Timer for adding tokens to the bucket
}

func NewTokenBucket(rate float64, capacity int) *TokenBucket {
	b := &TokenBucket{
		capacity: capacity,
		tokens:   capacity,
		rate:     rate,
		interval: time.Duration(float64(time.Second) / rate),
	}
	b.timer = time.AfterFunc(b.interval, b.refill)
	return b
}

func (b *TokenBucket) Take() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.tokens > 0 {
		b.tokens--
		return true
	}
	return false
}

func (b *TokenBucket) refill() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.tokens = b.capacity
	b.timer.Reset(b.interval)
}

func main() {
	bucket := NewTokenBucket(rateLimit, bucketSize)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if bucket.Take() {
			// Handle request...
			fmt.Fprintf(w, "Hello, World!")
		} else {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
		}
	})
	log.Fatal(http.ListenAndServe(":8000", nil))
}
