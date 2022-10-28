package main

import (
	"context"
	"log"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	var readers []chan int
	receiver := make(chan int)
	var wg = &sync.WaitGroup{}

	readers = startWorkers(readers, wg, ctx)
	split(ctx, receiver, readers)
	toReceiver(cancel, receiver)

	wg.Wait()
}

func toReceiver(cancel context.CancelFunc, receiver chan int) {
	go func() {
		for i := 0; i < 50; i++ {
			receiver <- i
		}

		cancel()
	}()
}

func split(ctx context.Context, receiver chan int, readers []chan int) {
	go func() {
		for {
			for _, reader := range readers {
				select {
				case <- ctx.Done():
					return
				case r := <- receiver:
					reader <- r
				}
			}
		}
	}()
}

func startWorkers(readers []chan int, wg *sync.WaitGroup, ctx context.Context) []chan int {
	for i := 0; i < 10; i++ {
		c := make(chan int)
		readers = append(readers, c)
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <- ctx.Done():
					return
				case r := <- c:
					log.Println("мидаеф 10 из 10: ", r)
				}
			}
		}()
	}

	return readers
}
