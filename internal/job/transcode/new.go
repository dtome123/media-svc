package transcode

import (
	"context"
	"log"
	"media-svc/internal/services"
	"media-svc/internal/services/media"
	"media-svc/internal/types"
	"sync"
	"time"
)

// Job holds information about a transcoding job
type Job struct {
	MediaID    string
	Status     string
	Error      string
	Result     transcodeResult
	CreatedAt  time.Time
	StartedAt  time.Time
	DoneAt     time.Time
	Renditions []types.Rendition
}

type transcodeResult struct {
	sourcePath string
	renditions []types.Rendition
}

// Orchestrator manages transcoding jobs and worker pool
type Orchestrator struct {
	svc        *services.Service              // service to handle media operations
	jobCh      chan media.TranscodeVideoInput // channel for job queue
	jobs       map[string]*Job                // map of jobs keyed by FilePath
	mu         sync.RWMutex                   // mutex to protect jobs map
	workers    int
	ctx        context.Context
	cancel     context.CancelFunc
	started    bool
	wg         sync.WaitGroup
	queueDepth int

	successChan chan *Job
	errorChan   chan *Job
}

// New creates an orchestrator with specified number of workers and queue depth
func New(svc *services.Service, workers int, queueDepth int) *Orchestrator {
	ctx, cancel := context.WithCancel(context.Background())
	return &Orchestrator{
		svc:        svc,
		jobCh:      make(chan media.TranscodeVideoInput, queueDepth),
		jobs:       make(map[string]*Job),
		workers:    workers,
		ctx:        ctx,
		cancel:     cancel,
		queueDepth: queueDepth,

		successChan: make(chan *Job, queueDepth),
		errorChan:   make(chan *Job, queueDepth),
	}
}

// Start launches the worker goroutines
func (o *Orchestrator) Start() {
	o.mu.Lock()
	if o.started {
		o.mu.Unlock()
		return
	}
	o.started = true
	o.mu.Unlock()

	for i := 0; i < o.workers; i++ {
		o.wg.Add(1)
		go o.worker(i)
	}
	log.Printf("Orchestrator started with %d workers", o.workers)
}

// Stop signals workers to stop and waits for them to finish
func (o *Orchestrator) Stop() {
	o.cancel()
	o.wg.Wait()
}

// AddJob adds a new transcoding job to the queue if not already present
func (o *Orchestrator) AddJob(in media.TranscodeVideoInput) {
	o.mu.Lock()
	if _, exists := o.jobs[in.MediaID]; !exists {
		o.jobs[in.MediaID] = &Job{
			MediaID:   in.MediaID,
			Status:    types.TranscodeJobStatusPending.String(),
			CreatedAt: time.Now(),
		}
	}
	o.mu.Unlock()

	// Blocking send to job channel; consider non-blocking or queue full logic if needed
	o.jobCh <- in
}

// GetJobStatus returns a copy of the job status or nil if job not found
func (o *Orchestrator) GetJobStatus(filePath string) *Job {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if job, ok := o.jobs[filePath]; ok {
		// Return a copy to avoid race conditions
		copyJob := *job
		return &copyJob
	}
	return nil
}

// worker is a goroutine that processes jobs from the job channel
func (o *Orchestrator) worker(id int) {
	defer o.wg.Done()
	log.Printf("worker-%d started", id)

	go o.processSuccess()
	go o.processError()
	go o.progressJob()

	for {
		select {
		case <-o.ctx.Done():
			log.Printf("worker-%d stopping (context cancelled)", id)
			return
		}

	}
}

func (o *Orchestrator) processSuccess() {
	for {
		select {
		case <-o.ctx.Done():
			return
		case msgs := <-o.successChan:
			o.svc.GetMediaSvc().UpdateTranscodeJobSuccess(context.Background(), media.UpdateTranscodeJobSuccessInput{
				MediaID:    msgs.MediaID,
				OutputPath: msgs.Result.sourcePath,
				Renditions: msgs.Result.renditions,
			})
		}
	}
}

func (o *Orchestrator) onSuccess(job *Job) {
	o.successChan <- job
}

func (o *Orchestrator) processError() {
	for {
		select {
		case <-o.ctx.Done():
			return
		case msgs := <-o.errorChan:
			o.svc.GetMediaSvc().UpdateTranscodeJobError(context.Background(), media.UpdateTranscodeJobErrorInput{
				MediaID: msgs.MediaID,
				Err:     msgs.Error,
			})
		}
	}
}

func (o *Orchestrator) onError(job *Job) {
	o.errorChan <- job
}

func (o *Orchestrator) progressJob() {
	for {
		select {
		case <-o.ctx.Done():
			return
		case input := <-o.jobCh:
			o.handleJob(input)
		}
	}
}

// handleJob executes the transcoding job by calling the media service
func (o *Orchestrator) handleJob(input media.TranscodeVideoInput) {
	// Mark job as processing
	o.mu.Lock()
	job := o.jobs[input.MediaID]
	if job == nil {
		job = &Job{
			MediaID:   input.MediaID,
			CreatedAt: time.Now(),
		}
		o.jobs[input.MediaID] = job
	}
	job.Status = types.TranscodeJobStatusProcessing.String()
	job.StartedAt = time.Now()
	o.mu.Unlock()

	// Call the TranscodeVideo service method

	result, err := o.svc.GetMediaSvc().TranscodeVideo(context.Background(), input)

	o.mu.Lock()
	defer o.mu.Unlock()

	if err != nil {
		job.Status = types.TranscodeJobStatusError.String()
		job.Error = err.Error()
		job.DoneAt = time.Now()
		o.onError(job)
		log.Printf("transcode failed for %s: %v", input.MediaID, err)
		return
	}

	job.Status = types.TranscodeJobStatusDone.String()
	job.DoneAt = time.Now()
	job.Result = transcodeResult{
		sourcePath: result.Path,
		renditions: result.Renditions,
	}
	o.onSuccess(job)
	log.Printf("transcode done for %s", input.MediaID)

}
