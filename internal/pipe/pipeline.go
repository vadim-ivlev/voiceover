package pipe

import (
	"fmt"
	"sync"
	"time"
)

// GenerateJobs - Stage1. creates jobs and sends them to the out channel.
// Parameters:
// numJobs: the number of jobs to create
// napTime: the time to sleep after creating a job
// out: a channel to send the jobs to
func GenerateJobs(numJobs, napTime int, out chan Job) {
	for i := 0; i < numJobs; i++ {
		job := createJob(i)
		LogJob(job, "created")
		out <- job
		if napTime > 0 {
			Nap(napTime)
		}
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
		logRecord := ProcessLogRecord{
			JobID:      job.ID,
			WorkerName: workerName,
			StartTime:  time.Now(),
		}

		job = operation(job)

		logRecord.EndTime = time.Now()
		logRecord.Duration = logRecord.EndTime.Sub(logRecord.StartTime)
		logRecord.DurationSeconds = logRecord.Duration.Seconds()

		job.ProcessLog = append(job.ProcessLog, logRecord)
		LogJob(job, "done by "+workerName)
		out <- job
	}
	// close(out)
}

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

func StartPipeline() {
	// Create a job queue
}
