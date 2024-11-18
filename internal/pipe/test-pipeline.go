package pipe

import (
	"sort"

	"github.com/rs/zerolog/log"
)

func StartTestPipeline() {
	// Create a job queue
	jobs := make(chan Job)
	go GenerateJobs(10, 0, jobs)

	squaredJobs := make(chan Job, 50)
	DoTeamWork(4, "SquareW", square, jobs, squaredJobs)

	cubedJobs := make(chan Job, 50)
	DoTeamWork(3, "CubeW", cube, squaredJobs, cubedJobs)

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
	// for i, job := range procssedJobs {
	// 	log.Info().Msgf("# %d: %s", i, job.String())
	// }
	log.Info().Msg(PrettyJSON(procssedJobs))
	// log.Info().Msg(CompactJSON(procssedJobs))

}
