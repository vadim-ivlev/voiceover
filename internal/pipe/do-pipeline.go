package pipe

// DoPipeline - starts the pipeline processing.
// Parameters:
//   - task: the task to process
//
// Returns:
//   - an array of processed jobs
//   - the number of jobs
//   - the base name of the output files
//   - an error if any
func DoPipeline(task Task) (doneJobs []Job, numJobs int, outputBaseName string, err error) {

	// get slice of epub text lines and the base name for the output files
	epubTextLines, outputBaseName, err := getEpubTextLines()
	if err != nil {
		return
	}

	numJobs = len(epubTextLines)

	// create an array of jobs
	jobs := newJobsArray(numJobs)
	// load the previous jobs to the new jobs array
	loadPreviuosJobs(jobs, task)

	// Create channels to pass the jobs between the workers
	jobsChan := make(chan Job, numJobs)
	textChan := make(chan Job, numJobs)
	translChan := make(chan Job, numJobs)
	// soundChan must be buffered because its jobs will be sorted after it is closed
	soundChan := make(chan Job, numJobs)

	// Fill the jobs channel with the jobs from a newly created array
	go toChannel(jobs, jobsChan)

	// add a text to each job and assign a voice
	// doTeamWork(4, "Text", "T", getTextOperation(epubTextLines), jobsChan, textChan)
	go DoWork(nil, "Add Text", "T", getTextOperation(epubTextLines), jobsChan, textChan)

	// translate the text.
	// go DoWork(nil, "Translate", "Tr", translateTextOperation, textChan, translChan)
	doTeamWork(30, "Translate", "Tr", translateTextOperation, textChan, translChan)

	// generate sound file for each job. Fan-out.
	doTeamWork(10, "Sound", "S", soundOperation, translChan, soundChan)

	// gatther the jobs into an array. Fan-in.
	doneJobs = toArray(soundChan)
	return doneJobs, numJobs, outputBaseName, nil
}
