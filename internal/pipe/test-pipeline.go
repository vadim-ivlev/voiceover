package pipe

import (
	"fmt"
	"sort"
	"sync"

	"github.com/rs/zerolog/log"
)

// DoTeamWork - creates a team of workers to process the jobs.
func DoTeamWork(workersNumber int, workerNamePrefix string, operation func(Job) Job, in <-chan Job, out chan<- Job) {
	// wait group to wait for the workers to finish
	wg := sync.WaitGroup{}
	for i := 0; i < workersNumber; i++ {
		wg.Add(1)
		workerName := fmt.Sprintf("%s-%d", workerNamePrefix, i)
		go DoWork(&wg, workerName, operation, in, out)
	}
	// wait for the workers to finish
	wg.Wait()
	close(out)
}

func StartTestPipeline() {
	// Create a job queue
	jobs := make(chan Job)
	go GenerateJobs(10, 0, jobs)

	squaredJobs := make(chan Job, 50)
	DoTeamWork(50, "SquareW", square, jobs, squaredJobs)

	cubedJobs := make(chan Job, 50)
	DoTeamWork(1, "CubeW", cube, squaredJobs, cubedJobs)

	// Consume the output
	log.Info().Msg("Waiting for the jobs to be processed >>>>>>>>>>>>>>>>>>>")
	procssedJobs := []Job{}
	for job := range cubedJobs {
		procssedJobs = append(procssedJobs, job)
	}

	// Sort the jobs
	sort.Slice(procssedJobs, func(i, j int) bool {
		return procssedJobs[i].ID < procssedJobs[j].ID
	})

	log.Info().Msg("All jobs have been processed.<<<<<<<<<<<<<<<<<<<")
	for i, job := range procssedJobs {
		log.Info().Msgf("# %d: %s", i, job.String())
	}

}
