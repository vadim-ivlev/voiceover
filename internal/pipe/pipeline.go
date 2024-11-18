package pipe

import (
	"sync"

	"github.com/rs/zerolog/log"
)

// GenerateJobs - Stage1. creates jobs and sends them to the out channel.
// Parameters:
// out: a channel to send the jobs to
func GenerateJobs(out chan Job) {
	for i := 0; i < 10; i++ {
		job := createJob(i)
		LogJob(job, "created")
		out <- job
		// Nap(1000)
	}
	close(out)
}

// DoWork - executes the operation on the job and sends the job to the out channel.
// Parameters:
// wg: a wait group to wait for the workers to finish
// workerName: the name of the worker
// operation: the operation to execute on the job
// in: a channel to receive the jobs from
// out: a channel to send the jobs to
func DoWork(wg *sync.WaitGroup, workerName string, operation func(Job) Job, in <-chan Job, out chan<- Job) {
	defer wg.Done()
	for job := range in {
		job = operation(job)
		job.ExecutedBy += workerName + ", "
		LogJob(job, "done by "+workerName)
		out <- job
	}
	// close(out)
}

func StartPipeline() {
	// Create a job queue
	jobs := make(chan Job, 50)
	go GenerateJobs(jobs)

	squaredJobs := make(chan Job, 50)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go DoWork(&wg, "SquareW1", square, jobs, squaredJobs)
	wg.Add(1)
	go DoWork(&wg, "SquareW2", square, jobs, squaredJobs)
	wg.Add(1)
	go DoWork(&wg, "SquareW3", square, jobs, squaredJobs)

	// wait for the workers to finish
	wg.Wait()
	close(squaredJobs)

	cubedJobs := make(chan Job, 50)
	wg.Add(1)
	go DoWork(&wg, "CubeW1", cube, squaredJobs, cubedJobs)
	wg.Wait()
	close(cubedJobs)

	// Consume the output
	log.Info().Msg("Waiting for the jobs to be processed............................")
	procssedJobs := []Job{}
	for job := range cubedJobs {
		procssedJobs = append(procssedJobs, job)
	}

	log.Info().Msg("All jobs have been processed.")
	for i, job := range procssedJobs {
		log.Info().Msgf("# %d: %s", i, job.String())
	}

}
