package tworker

import (
	"fmt"

	"github.com/cihub/seelog"
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

func (w *WorkerTb) Init(workerID int, pool *WorkerPool) {
	w.pool = pool
	w.workerID = workerID
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

WORK_LOOP:
	for {
		select {
		case <-w.pool.doneContext.Done():
			{
				break WORK_LOOP
			}
		}
	}

	w.linfo("Done ...")
}
