package tworker

type IWorker interface {
	Init(int, *WorkerPool) error
	Run()
	New() IWorker
}
