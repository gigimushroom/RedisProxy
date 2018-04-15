package util

// Scheduler schedule a job to an available worker from the worker pool
type Scheduler struct {
	WorkerPool chan chan Job
	maxWorkerNum int
	JobQueue chan Job
	workers []Worker
}

// NewScheduler create a new Scheduler
func NewScheduler(maxWorkers int, jobMax int) *Scheduler {
	pool := make(chan chan Job, maxWorkers)

	jobQueue := make(chan Job, jobMax)
	workers := make([]Worker, maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		worker := NewWorker(pool)
		workers = append(workers, worker)
	}
	return &Scheduler{pool, maxWorkers, jobQueue, workers}
}

// Run starts workers and kick the schedule routine
func (d *Scheduler) Run() {
	for _,w := range d.workers {
		w.Start()
	}
	go d.schedule()
}

func (d *Scheduler) schedule() {
	for {
		select {
		case job := <-d.JobQueue:
			// After receives a job, find an available worker
			// Send worker the job
			go func(job Job) {
				jobChannel := <-d.WorkerPool
				jobChannel <- job
			}(job)
		}
	}
}

func (d *Scheduler) Stop() {
	for _,w := range d.workers {
		w.Stop()
	}
}

