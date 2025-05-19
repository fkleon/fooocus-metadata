package fooocusplus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePrivateLog(t *testing.T) {
	// Private log with mix of image types
	var privateLogFile = "./testdata/log.html"
	images, err := ParsePrivateLog(privateLogFile)
	require.NoError(t, err)

	var expected = []string{
		"fooocusplus-meta.png",
		"fooocusplus-meta.jpeg",
		"fooocusplus-meta.webp",
		"2025-04-23_11-27-25_6011.png",
	}
	assert.Len(t, images, len(expected))

	for _, image := range expected {
		t.Run(image, func(t *testing.T) {
			assert.Contains(t, images, image)

			metadata := images[image]
			require.NotZero(t, metadata)
			require.IsType(t, Metadata{}, metadata)

			assert.Equal(t, "FooocusPlus 1.0.0", metadata.Version)
			assert.Equal(t, "Fooocus", metadata.MetadataScheme)
			// TODO reference data to compare to
			//assert.Equal(t, meta, images[image])
		})
	}
}
