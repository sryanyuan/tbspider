package tworker

type IWorker interface {
	Init(int, *WorkerPool)
	Run()
	New() IWorker
}
