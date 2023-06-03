package utils_test

import (
	"github.com/celsopires1999/encoder/framework/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsJson(t *testing.T) {
	t.Run("Should return true if json is valid", func(t *testing.T) {
		json := `{
			"id": "525b5fd9-700d-4feb-89c0-415a1e6e148c",
			"file_path": "convite.mp4",
			"status": "pending"
		  }`

		err := utils.IsJson(json)
		require.Nil(t, err)
	})

	t.Run("Should return false if json is invalid", func(t *testing.T) {
		json := `wes`
		err := utils.IsJson(json)
		require.Error(t, err)
	})
}
