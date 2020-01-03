package concurrence

import (
	"fmt"
	"sync"
)

// --------------------------- Job ---------------------
type Job interface {
	Do()
}

// --------------------------- Worker ---------------------
type Worker struct {
	JobQueue chan Job
	closed   chan struct{}
	wg       *sync.WaitGroup
}

func NewWorker() *Worker {
	return &Worker{
		JobQueue: make(chan Job),
		closed:   make(chan struct{}),
	}
}
func (w *Worker) Run(wq chan chan Job) {
	go func(wg *sync.WaitGroup) {
		defer func() {
			if w.wg != nil {
				w.wg.Done()
			}
		}()
		for {
			//把worker得工作队列放到wq中
			wq <- w.JobQueue
			select {
			case job := <-w.JobQueue:
				//如果工作队列分配到一个job
				job.Do()
			case <-w.closed:
				//如果收到关闭信号
				fmt.Println("worker stopped")
				return
			}
		}
	}(w.wg)
}

func (w *Worker) Stop(wg *sync.WaitGroup) {
	w.wg = wg
	close(w.closed)
}
