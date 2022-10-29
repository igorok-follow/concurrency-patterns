package main

import (
	"context"
	"log"
	"strconv"
	"sync"
)

type Processor struct {
	Receiver   chan *Job
	Readers    []chan *Job
	ReadersNum int
	JobsNum    int
	Wg         *sync.WaitGroup
}

type Job struct {
	Id int
}

func newProcessor(readersNum int, jobsNum int) *Processor {
	return &Processor{
		Receiver:   make(chan *Job),
		Readers:    nil,
		ReadersNum: readersNum,
		JobsNum:    jobsNum,
		Wg:         new(sync.WaitGroup),
	}
}

func (p *Processor) run(ctx context.Context, cancel context.CancelFunc) {
	p.initWorkers(ctx)
	p.split(ctx)
	p.toReceiver(cancel)
}

func (p *Processor) initWorkers(ctx context.Context) {
	for i := 0; i < p.ReadersNum; i++ {
		c := make(chan *Job)
		p.Readers = append(p.Readers, c)

		p.Wg.Add(1)
		go func() {
			defer p.Wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case j := <-c:
					log.Println("JOB " + strconv.Itoa(j.Id) + " FINISHED")
				}
			}
		}()
	}
}

func (p *Processor) split(ctx context.Context) {
	go func() {
		for {
			for _, reader := range p.Readers {
				select {
				case <-ctx.Done():
					return
				case j := <-p.Receiver:
					reader <- j
				}
			}
		}
	}()
}

func (p *Processor) toReceiver(cancel context.CancelFunc) {
	go func() {
		for i := 0; i < p.JobsNum; i++ {
			p.Receiver <- &Job{Id: i}
		}

		cancel()
	}()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	processor := newProcessor(10, 100)
	processor.run(ctx, cancel)

	processor.Wg.Wait()
}
