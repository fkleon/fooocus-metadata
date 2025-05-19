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

	testCases := []struct {
		file             string
		expectedScheme   string
		expectedSoftware string
	}{
		{"fooocus-meta.png", "fooocus", "Fooocus v2.5.5"},
		{"fooocus-meta.jpeg", "fooocus", "Fooocus v2.5.5"},
		{"fooocus-meta.webp", "fooocus", "Fooocus v2.5.5"},
		{"a1111-meta.png", "a1111", "Fooocus v2.5.5"},
		{"a1111-meta.jpeg", "a1111", "Fooocus v2.5.5"},
		{"a1111-meta.webp", "a1111", "Fooocus v2.5.5"},
		{"2024-01-05_23-11-48_9167.png", "fooocus", "v2.1.860"},
	}

	assert.Len(t, images, len(testCases))

	for _, tc := range testCases {
		t.Run(tc.file, func(t *testing.T) {
			assert.Contains(t, images, tc.file)

			metadata := images[tc.file]
			require.NotZero(t, metadata)
			require.IsType(t, Metadata{}, metadata)

			assert.Equal(t, tc.expectedSoftware, metadata.Version)
			assert.Equal(t, tc.expectedScheme, metadata.MetadataScheme)
			// TODO reference data to compare to
			//assert.Equal(t, meta, images[image])
		})
	}
}
