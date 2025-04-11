package fooocus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePrivateLog(t *testing.T) {
	// Private log with mix of legacy and current metadata format
	var privateLogFile = "./testdata/log.html"
	images, err := ParsePrivateLog(privateLogFile)
	require.NoError(t, err)

	var expected = []string{
		"fooocus-meta.png",
		"fooocus-meta.jpeg",
		"fooocus-meta.webp",
		"a111-meta.png",
		"a111-meta.jpeg",
		"a111-meta.png",
		"2024-01-05_23-11-48_9167.png",
	}
	assert.Len(t, images, len(expected))

	for _, image := range expected {
		assert.Contains(t, images, image)
		assert.NotNil(t, images[image])
	}
}
