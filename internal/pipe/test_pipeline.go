package pipe

import (
	"sync"

	"github.com/rs/zerolog/log"
)

func StartTestPipeline() {
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
	wg1 := sync.WaitGroup{}
	wg1.Add(1)
	go DoWork(&wg1, "CubeW1", cube, squaredJobs, cubedJobs)
	wg1.Wait()
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
