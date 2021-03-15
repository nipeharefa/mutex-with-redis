package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
)

// var ctx = context.Background()
// u[dateCart with backoff options
func updateCart(locker *redislock.Client, c chan int, ctx context.Context, i int) {

	backoff := redislock.LinearBackoff(10 * time.Millisecond)
	lock, err := locker.Obtain(ctx, "my-key", 2*time.Second, &redislock.Options{
		RetryStrategy: backoff,
	})
	if err == redislock.ErrNotObtained {
		fmt.Println("Could not obtain lock!", lock)
	} else if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("lock saya", i)
	lock.Release(ctx)
	fmt.Println("release", i)
	c <- 1
}

func sendUpdateCart(locker *redislock.Client) {
	//
	v := 10
	ctx := context.Background()
	jobs := make(chan int, v)

	for i := 0; i < v; i++ {
		go updateCart(locker, jobs, ctx, i)
	}

	for {
		<-jobs
		fmt.Println("next")
	}

}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	defer rdb.Close()

	// Create a new lock client.
	locker := redislock.New(rdb)

	sendUpdateCart(locker)

}
