package main

import (
	"log"
	"strconv"
	"time"
)

type Processor struct {
	Jobs    chan *Job
	Done    chan *Worker
	Workers []*Worker
}

type Job struct {
	Id   int
	Name string
}

type Worker struct {
	Name string
}

func NewProcessor(workersNum int) *Processor {
	workers := make([]*Worker, workersNum)
	for i := 0; i < len(workers); i++ {
		workers[i] = &Worker{
			Name: "Worker Num " + strconv.Itoa(i),
		}
	}

	return &Processor{
		Jobs:    make(chan *Job),
		Done:    make(chan *Worker),
		Workers: workers,
	}
}

func (p *Processor) Run() {
	go func() {
		for {
			select {
			default:
				if len(p.Workers) > 0 {
					w := p.Workers[0]
					p.Workers = p.Workers[1:]
					w.RunJob(<-p.Jobs, p.Done)
				}
			case w := <-p.Done:
				p.Workers = append(p.Workers, w)
			}
		}
	}()
}

func GetJob(id int) <-chan *Job {
	c := make(chan *Job)

	go func() {
		for {
			job := &Job{
				Id: id,
			}

			c <- job
		}
	}()

	return c
}

func (w *Worker) RunJob(job *Job, done chan *Worker) {
	go func() {
		time.Sleep(time.Second * 1)
		log.Println("JOB " + strconv.Itoa(job.Id) + " FINISHED")
		done <- w
	}()
}

func (p *Processor) ScheduleJob(job <-chan *Job) {
	j := <-job
	p.Jobs <- j
}

func main() {
	processor := NewProcessor(5)
	processor.Run()

	for i := 0; i < 100; i++ {
		processor.ScheduleJob(GetJob(i))
	}
	time.Sleep(time.Second * 5)
}
