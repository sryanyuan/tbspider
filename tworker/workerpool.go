package tworker

import (
	"fmt"
	"sync"

	"github.com/sryanyuan/tbspider/tconfig"

	"golang.org/x/net/context"
)

// for register
var (
	registeredWorkers map[string]IWorker
)

func registerWorker(workerName string, worker IWorker) {
	if nil == registeredWorkers {
		registeredWorkers = make(map[string]IWorker)
	}
	registeredWorkers[workerName] = worker
}

func createWorker(workerName string) IWorker {
	v, ok := registeredWorkers[workerName]
	if !ok {
		return nil
	}

	return v.New()
}

type WorkerPool struct {
	workers     []IWorker
	fnCancel    context.CancelFunc
	doneContext context.Context
	wg          sync.WaitGroup
}

func NewWorkerPool() *WorkerPool {
	ins := &WorkerPool{
		workers: make([]IWorker, 0, 32),
	}
	return ins
}

func (w *WorkerPool) InitWithWorkerCount(count int) error {
	if nil == registeredWorkers {
		return fmt.Errorf("no worker registered")
	}
	config := tconfig.StoreConfig(nil)
	w.workers = make([]IWorker, count)

	// initialize context
	ctx := context.Background()
	w.doneContext, w.fnCancel = context.WithCancel(ctx)

	for i, _ := range w.workers {
		worker := createWorker(config.WorkerName)
		if nil == worker {
			return fmt.Errorf("failed to create worker [%s]", config.WorkerName)
		}
		worker.Init(i, w)
		w.workers[i] = worker
	}

	// run all workers
	for _, v := range w.workers {
		w.wg.Add(1)
		go v.Run()
	}

	return nil
}

func (w *WorkerPool) WaitWorkersDone() {
	w.fnCancel()
	w.wg.Wait()
}
