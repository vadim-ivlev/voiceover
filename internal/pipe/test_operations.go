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

// square - square the number
func square(job Job) Job {
	Nap(2000)
	job.Test += fmt.Sprintf("Job %d squared. ", job.ID)
	return job
}

// cube - cube the number
func cube(job Job) Job {
	Nap(100)
	job.Test += fmt.Sprintf("Job %d cubed. ", job.ID)
	return job
}
