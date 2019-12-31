package concurrence

import "fmt"

// --------------------------- Job ---------------------
type Job interface {
	Do()
}

// --------------------------- Worker ---------------------
type Worker struct {
	JobQueue chan Job
	closed   chan struct{}
}

func NewWorker() Worker {
	return Worker{
		JobQueue: make(chan Job),
		closed:   make(chan struct{}),
	}
}
func (w Worker) Run(wq chan chan Job) {
	go func() {
		for {
			wq <- w.JobQueue
			select {
			case job := <-w.JobQueue:
				job.Do()
			case <-w.closed:
				fmt.Println("worker stopped")
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	close(w.closed)
}
