package util

// JobHandler contains the actual work belong to a job
type JobHandler interface {
	Lookup(string) string
}

type Job struct{
	Request string
	JobHandler JobHandler
	Resp chan string
}

type Worker struct {
	WorkerPool  chan chan Job
	JobChannel  chan Job
	quit    	chan bool
}

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
}

func (w Worker) Start() {
	go func() {
		for {
			// Register myself to worker pool
			w.WorkerPool <- w.JobChannel
			select {
			case job := <-w.JobChannel:
				// Do the job
				val := job.JobHandler.Lookup(job.Request)
				// Send response back to caller
				job.Resp <- val

			case <-w.quit:
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}