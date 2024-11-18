package pipe

import (
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/exp/rand"
)

// Nap - Sleep for a random duration between 0 and milliseconds
func Nap(milliseconds int) {
	// random duration between 0 and 1 second
	sleepDuration := time.Duration(rand.Intn(milliseconds)) * time.Millisecond
	time.Sleep(sleepDuration)
}

// square - square the number
func square(job Job) Job {
	Nap(2000)
	return job
}

// cube - cube the number
func cube(job Job) Job {
	Nap(100)
	return job
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
