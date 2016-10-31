package tworker

import (
	"fmt"

	"github.com/cihub/seelog"
	"github.com/sryanyuan/tbspider/tconfig"
)

// shared work task
type tbWorkTask struct {
	url string
}

var (
	sharedTbTaskQueue []*tbWorkTask
)

type WorkerTb struct {
	pool     *WorkerPool
	workerID int
}

func init() {
	registerWorker("tumblr", &WorkerTb{})
}

func (w *WorkerTb) linfo(args ...interface{}) {
	l := fmt.Sprintln(args...)
	seelog.Info("WorkerTb[", w.workerID, "] : ", l)
}

func (w *WorkerTb) Init(workerID int, pool *WorkerPool) error {
	w.pool = pool
	w.workerID = workerID

	// here we initialize work task once
	if nil == sharedTaskQueue {
		seelog.Info("Initialize task queue, it may takes some time, please wait ...")
	}

	return nil
}

func (w *WorkerTb) New() IWorker {
	n := &WorkerTb{}
	return n
}

func (w *WorkerTb) Run() {
	w.linfo("Running ...")
	defer func() {
		w.pool.wg.Done()
	}()

	// we get the worker count
	workerCount := tconfig.StoreConfig(nil).MaxWorkers
	totalTaskCount := len(sharedTbTaskQueue)
	currentTaskCount := totalTaskCount / workerCount
	currentWorkingStartIndex := w.workerID * currentTaskCount
	reminderTaskCount := totalTaskCount % workerCount
	if w.workerID == workerCount-1 {
		// last one, should do the left work
		currentTaskCount += reminderTaskCount
	}

	for taskIndex := currentWorkingStartIndex; taskIndex < currentWorkingStartIndex+currentTaskCount; taskIndex++ {
		task := sharedTbTaskQueue[taskIndex]

		// do get
	}

	w.linfo("Done ...")
}
