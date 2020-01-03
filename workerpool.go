package concurrence

import "fmt"

// --------------------------- WorkerPool ---------------------
type WorkerPool struct {
	workerlen   int
	JobQueue    chan Job
	WorkerQueue chan chan Job
}

func NewWorkerPool(workerlen int) *WorkerPool {
	return &WorkerPool{
		workerlen:   workerlen,
		JobQueue:    make(chan Job),
		WorkerQueue: make(chan chan Job, workerlen),
	}
}
func (wp *WorkerPool) Run() (workers []*Worker) {
	fmt.Println("初始化worker")
	//初始化worker
	for i := 0; i < wp.workerlen; i++ {
		worker := NewWorker()
		worker.Run(wp.WorkerQueue)
		workers = append(workers, worker)
	}
	// 循环获取可用的worker,往worker中写job
	go func() {
		for {
			select {
			case job := <-wp.JobQueue:
				//从pool里面找到一个chan Job（属于某个worker）
				jobCh := <-wp.WorkerQueue
				//把当前得job分配给这个chan Job（属于某个worker）
				jobCh <- job
			}
		}
	}()
	return
}
