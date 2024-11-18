package pipe

import (
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
	return job
}

// cube - cube the number
func cube(job Job) Job {
	Nap(100)
	return job
}
