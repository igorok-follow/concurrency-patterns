package main

import "context"

type Processor struct {
	Ctx context.Context
}

type Job struct {
	Id int
}

func split(receiver chan *Job) chan *Job {
	go func() {
		for {
			select {
			case <- ctx.Done():
				return
			case j := <-
			}
		}
	}()
}

func initWorkers(ctx context.Context, readers chan *Job) chan *Job {

}

func toReceiver(receiver chan *Job) {

}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	receiver := make(chan *Job)
	var readers chan *Job

	readers = initWorkers(ctx)
	split(receiver)
	toReceiver(receiver)
}