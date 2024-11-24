package pipe

import (
	"github.com/rs/zerolog/log"
)

// square - square the number
func square(job Job) (Job, error) {
	Nap()
	return job, nil
}

// cube - cube the number
func cube(job Job) (Job, error) {
	Nap()
	return job, nil
}

func StartTestPipeline() {
	// Create a job queue
	jobs := make(chan Job)
	go GenerateJobs(10, 0, jobs)

	squaredJobs := make(chan Job, 50)
	doTeamWork(4, "SquareW", square, jobs, squaredJobs)

	cubedJobs := make(chan Job, 50)
	doTeamWork(3, "CubeW", cube, squaredJobs, cubedJobs)

	// gatther the jobs into an array
	processedJobs := toArray(cubedJobs)
	log.Info().Msg(PrettyJSON(processedJobs))
}
