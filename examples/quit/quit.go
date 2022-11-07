package main

import (
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	quit := make(chan string)

	c := initChan("data", quit)
	for i := rand.Intn(10); i != 0; i-- {
		log.Println(i)
		log.Println(<-c)
	}
	quit <- "End!"

	log.Println(<-quit)
}

func initChan(data string, quit chan string) chan string {
	c := make(chan string)

	go func() {
		for i := 0; i < 10; i++ {
			select {
			case c <- data:
				// do nothing ...
			case <-quit:
				quit <- "Bye!"
				return
			}
		}
	}()

	return c
}
