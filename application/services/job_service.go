package services

import (
	"errors"
	"os"
	"strconv"

	"github.com/celsopires1999/encoder/application/repositories"
	"github.com/celsopires1999/encoder/domain"
)

type JobService struct {
	Job           *domain.Job
	JobRepository repositories.JobRepository
	VideoService  VideoService
}

func (j *JobService) Start() error {
	if err := j.changeJobStatus("DOWNLOADING"); err != nil {
		return j.failJob(err)
	}
	if err := j.VideoService.Download(os.Getenv("inputBucketName")); err != nil {
		return j.failJob(err)
	}

	if err := j.changeJobStatus("FRAGMENTING"); err != nil {
		return j.failJob(err)
	}
	if err := j.VideoService.Fragment(); err != nil {
		return j.failJob(err)
	}

	if err := j.changeJobStatus("ENCODING"); err != nil {
		return j.failJob(err)
	}
	if err := j.VideoService.Encode(); err != nil {
		return j.failJob(err)
	}

	if err := j.performUpload(); err != nil {
		return j.failJob(err)
	}

	if err := j.changeJobStatus("COMPLETED"); err != nil {
		return j.failJob(err)
	}

	return nil
}

func (j *JobService) performUpload() error {
	if err := j.changeJobStatus("UPLOADING"); err != nil {
		return j.failJob(err)
	}

	videoUpload := NewVideoUpload()
	videoUpload.OutputBucket = os.Getenv("outputBucketName")
	videoUpload.VideoPath = os.Getenv("localStoragePath") + "/" + j.VideoService.Video.ID

	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY_UPLOAD"))
	if err != nil {
		return err
	}

	doneUplod := make(chan string)

	go videoUpload.ProcessUpload(concurrency, doneUplod)

	uploadResult := <-doneUplod

	if uploadResult != "upload completed" {
		return j.failJob(errors.New(uploadResult))
	}

	return nil
}

func (j *JobService) changeJobStatus(status string) error {
	var err error

	j.Job.Status = status
	j.Job, err = j.JobRepository.Update(j.Job)
	if err != nil {
		return j.failJob(err)
	}

	return err
}

func (j *JobService) failJob(error error) error {
	j.Job.Status = "FAILED"
	j.Job.Error = error.Error()

	if _, err := j.JobRepository.Update(j.Job); err != nil {
		return err
	}

	return error
}
