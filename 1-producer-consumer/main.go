//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Ci sono due approcci: un canale bool per quando finisce il ciclo for di consumer e uno con un wg che attende la fine dell'esecuzione delle goroutines

func producer(stream Stream, tweets chan *Tweet, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		tweet, err := stream.Next()
		if errors.Is(err, ErrEOF) {
			close(tweets)
			return
		}
		// Al canale tweets passo i tweet
		tweets <- tweet
	}
}

func consumer(tweets chan *Tweet, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	tweets := make(chan *Tweet)
	var wg sync.WaitGroup
	wg.Add(2)
	// Producer
	go producer(stream, tweets, &wg)

	// Consumer
	go consumer(tweets, &wg)

	wg.Wait()

	fmt.Printf("Process took %s\n", time.Since(start))
}
