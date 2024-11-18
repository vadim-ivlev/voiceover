package pipe

import (
	"github.com/rs/zerolog/log"
)

// LogJob - logs the job information
func LogJob(job Job, msg string) {
	log.Info().Msgf("Job %d %s", job.ID, msg)
	// log.Info().Msg(job.String())
}

// GenerateJobs - Stage1. creates jobs and sends them to the out channel.
// Parameters:
// out: a channel to send the jobs to
func GenerateJobs(out chan Job) {
	for i := 0; i < 5; i++ {
		job := createJob(i)
		LogJob(job, "created")
		out <- job
		Nap(1000)
	}
	close(out)
}

// DoWork - executes the operation on the job and sends the job to the out channel.
// Parameters:
// workerName: the name of the worker
// operation: the operation to execute on the job
// in: a channel to receive the jobs from
// out: a channel to send the jobs to
func DoWork(workerName string, operation func(Job) Job, in <-chan Job, out chan<- Job) {
	for job := range in {
		job = operation(job)
		job.ExecutedBy += workerName + ", "
		LogJob(job, "done by "+workerName)
		out <- job
	}
	close(out)
}

func StartPipeline() {
	// Create a job queue
	jobs := make(chan Job, 50)
	go GenerateJobs(jobs)

	squaredJobs := make(chan Job, 50)
	go DoWork("SqureWorker1", square, jobs, squaredJobs)

	cubedJobs := make(chan Job, 50)
	go DoWork("CubeWorker1", cube, squaredJobs, cubedJobs)

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
