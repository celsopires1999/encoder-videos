package domain_test

import (
	"testing"
	"time"

	"github.com/celsopires1999/encoder/domain"
	"github.com/stretchr/testify/require"
)

func TestValidateErrorWhenVideoIsEmpty(t *testing.T) {
	expectedError := "encoded_video_folder:  does not validate as uuid;file_path:  does not validate as notnull;resource_id:  does not validate as notnull"
	video := domain.NewVideo()
	err := video.Validate()
	require.Error(t, err)
	require.Equal(t, expectedError, err.Error())
}

func TestValidateErrorWhenIdIsInvalid(t *testing.T) {
	expectedError := "encoded_video_folder: fake-id does not validate as uuid"
	video := domain.NewVideo()
	video.ID = "fake-id"
	video.ResourceID = "resource-id"
	video.FilePath = "file-path"
	video.CreatedAt = time.Now()
	err := video.Validate()
	require.Error(t, err)
	require.Equal(t, expectedError, err.Error())
}
