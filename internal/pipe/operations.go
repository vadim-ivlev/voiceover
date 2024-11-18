package pipe

import "github.com/rs/zerolog/log"

// createJob - create a job with the given id
func createJob(id int) Job {
	job := Job{
		ID: id,
	}
	return job
}

// LogJob - logs the job information
func LogJob(job Job, msg string) {
	log.Info().Msgf("Job %d %s", job.ID, msg)
	// log.Info().Msg(job.String())
}
