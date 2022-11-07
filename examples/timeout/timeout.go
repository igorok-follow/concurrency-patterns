package main

import (
	"log"
	"time"
)

func main() {
	c := initChan("data")

	select {
	case d := <-c:
		log.Println(d)
	case <-time.After(time.Second):
		log.Println("you are too slow...")
	}
}

func initChan(data string) <-chan string {
	c := make(chan string)

	go func() {
		time.Sleep(time.Millisecond * 1001)
		c <- data
	}()

	return c
}
