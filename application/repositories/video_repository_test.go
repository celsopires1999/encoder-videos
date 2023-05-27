package repositories_test

import (
	"testing"
	"time"

	"github.com/celsopires1999/encoder/application/repositories"
	"github.com/celsopires1999/encoder/domain"
	"github.com/celsopires1999/encoder/framework/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestVideoRepositoryDbInsert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.New().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()
	video.ResourceID = "fake-resource-id"

	repo := repositories.VideoRepositoryDb{Db: db}

	v, err := repo.Insert(video)
	require.Nil(t, err)
	require.NotNil(t, v)

	v, err = repo.Find(video.ID)
	require.Nil(t, err)
	require.NotNil(t, v)
	require.Equal(t, v.ID, video.ID)
	require.Equal(t, v.ResourceID, video.ResourceID)
	require.Equal(t, v.FilePath, video.FilePath)
	require.Equal(t, v.Jobs, video.Jobs)
	require.Equal(t, v.CreatedAt.Local(), video.CreatedAt.Local())
}
