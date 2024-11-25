package pipe

import (
	"fmt"
	"sync"
	"time"
)

// newJobsArray - creates an array of jobs.
func newJobsArray(numJobs int) []Job {
	jobs := make([]Job, numJobs)
	for i := 0; i < numJobs; i++ {
		jobs[i] = Job{
			ID: i,
		}
	}
	return jobs
}

// toChannel - moves the jobs from the array to the channel
// Parameters:
// jobsArray: the array of jobs
// jobsChan: the channel to send the jobs to
func toChannel(jobsArray []Job, jobsChan chan Job) {
	for _, job := range jobsArray {
		jobsChan <- job
		Nap()
	}
	close(jobsChan)
}

// DoWork - executes the operation on the job and sends the job to the out channel.
// Parameters:
// wg: a wait group to wait for the workers to finish. Can be nil if the workers works alone, outside doTeamWork.
// work: the name of the work
// workerName: the name of the worker
// operation: the operation to execute on the job
// in: a channel to receive the jobs from
// out: a channel to send the jobs to
func DoWork(wg *sync.WaitGroup, work, workerName string, operation JobFunction, in <-chan Job, out chan<- Job) {
	// Decrease the wait group counter when the worker belongs to a team.
	if wg != nil {
		defer wg.Done()
	}
	for job := range in {
		// Check if work needs to be done
		if ShouldSkipWork(job, work, workerName) {
			out <- job
			continue
		}

		logRecord := ProcessLogRecord{
			JobID:     job.ID,
			Work:      work,
			Worker:    workerName,
			StartTime: time.Now(),
		}
		Nap()
		job, err := operation(job)
		if err != nil {
			logRecord.Error = err.Error()
		}
		logRecord.EndTime = time.Now()
		logRecord.Duration = logRecord.EndTime.Sub(logRecord.StartTime)
		logRecord.DurationSeconds = logRecord.Duration.Seconds()

		job.ProcessLog = append(job.ProcessLog, logRecord)
		LogJob(job, fmt.Sprintf("done by %10s in %.3f seconds", workerName, logRecord.DurationSeconds))
		out <- job
	}
	// Close channel if worker works alone.
	if wg == nil {
		close(out)
	}
}

// doTeamWork - creates a team of workers to process the jobs.
// Parameters:
// workersNumber: the number of workers in the team
// work: the name of the work
// workerNamePrefix: the prefix for the worker names
// operation: the operation to execute on the job
// in: a channel to receive the jobs from
// out: a channel to send the jobs to
func doTeamWork(workersNumber int, work, workerNamePrefix string, operation JobFunction, in <-chan Job, out chan<- Job) {
	// wait group to wait for the workers to finish
	wg := sync.WaitGroup{}
	for i := 0; i < workersNumber; i++ {
		wg.Add(1)
		workerName := fmt.Sprintf("%s-%d", workerNamePrefix, i)
		go DoWork(&wg, work, workerName, operation, in, out)
	}
	// wait for the workers to finish
	wg.Wait()
	close(out)
}

// ShouldSkipWork - checks if the work should be skipped.
func ShouldSkipWork(job Job, work, workerName string) bool {
	// If no process log, the work should not be skipped
	if job.ProcessLog == nil || len(job.ProcessLog) == 0 {
		return false
	}

	// check if this work has already been done
	for _, logRecord := range job.ProcessLog {
		if logRecord.Work == work {
			LogJob(job, fmt.Sprintf("Work %s skipped by %s because it has already been done", work, workerName))
			return true
		}
	}

	// check if the previous worker has finished the job successfully
	lastLog := job.ProcessLog[len(job.ProcessLog)-1]
	if lastLog.Error != "" {
		LogJob(job, fmt.Sprintf("Work %s skipped by %s because the previous worker has failed", work, workerName))
		return true
	}

	return false

}
