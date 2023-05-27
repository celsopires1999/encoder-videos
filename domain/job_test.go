package domain_test

import (
	"testing"
	"time"

	"github.com/celsopires1999/encoder/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewJob(t *testing.T) {
	video := domain.NewVideo()
	video.ID = uuid.New().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	job, err := domain.NewJob("path", "converted", video)
	assert.NotNil(t, job)
	assert.Nil(t, err)
	assert.Nil(t, func(id string) error {
		_, err := uuid.Parse(id)
		return err
	}(job.ID))
}
