package main

import (
	"fmt"
	"github.com/ghjan/concurence"
	"os"
	"os/signal"
	"runtime"
	"time"
)

//定义一个实现Job接口的数据
type Score struct {
	Num int
}

var Workers []concurrence.Worker

//定义对数据的处理
func (s *Score) Do() {
	time.Sleep(1 * time.Second)
	fmt.Println("num:", s.Num)
}

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	printInfo()
	//// 定义一个cron运行器
	//cr := cron.New()
	//// 定时5秒，每1秒执行handle
	//cr.AddFunc("*/1 * * * * *", handleNumGoroutine)
	//
	//// 开始
	//cr.Start()
	//defer cr.Stop()

	num := 100 * 100 * 20
	// debug.SetMaxThreads(num + 1000) //设置最大线程数
	// 注册工作池，传入任务
	// 参数1 worker并发个数
	p := concurrence.NewWorkerPool(num)
	Workers = p.Run()

	//写入datanum条数据
	dataNum := 100 * 100 * 100
	go func() {
		for i := 1; i <= dataNum; i++ {
			sc := &Score{Num: i}
			p.JobQueue <- sc //数据传进去会被自动执行Do()方法，具体对数据的处理自己在Do()方法中定义
		}
	}()

	select {
	case sig := <-c:
		fmt.Printf("Got %s signal. Aborting...\n", sig)
		for _, worker := range Workers {
			worker.Stop()
		}
	}
}

func printInfo() {
	//定时打印输出当前进程的Goroutine 个数
	//ticker := time.NewTicker(time.Second * 1)
	t1 := time.NewTimer(time.Second * 1)
	go func() {
		for {
			select {
			case <-t1.C:
				t1.Reset(time.Second * 1)
				handleNumGoroutine()
			}
		}
	}()
}

func handleNumGoroutine() {
	fmt.Println("runtime.NumGoroutine() :", runtime.NumGoroutine())
}
