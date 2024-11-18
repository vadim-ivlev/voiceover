package pipe

import (
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

// Nap - Sleep for a random duration between 0 and milliseconds
func Nap(milliseconds int) {
	// random duration between 0 and 1 second
	sleepDuration := time.Duration(rand.Intn(milliseconds)) * time.Millisecond
	time.Sleep(sleepDuration)
}

func createJob(id int) Job {
	job := Job{
		ID:    id,
		JobID: fmt.Sprintf("Job%08d", id),
	}
	return job
}

// square - square the number
func square(job Job) Job {
	Nap(1000)
	job.Test += fmt.Sprintf("Job %d squared. ", job.ID)
	return job
}

// cube - cube the number
func cube(job Job) Job {
	Nap(1000)
	job.Test += fmt.Sprintf("Job %d cubed. ", job.ID)
	return job
}
