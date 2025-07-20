package fooocus

import (
	"maps"
	"slices"
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
		file     string
		scheme   string
		software string
		meta     *MetadataV23
	}{
		{"fooocus-meta.png", "fooocus", "Fooocus v2.5.5", nil},
		{"fooocus-meta.jpeg", "fooocus", "Fooocus v2.5.5", nil},
		{"fooocus-meta.webp", "fooocus", "Fooocus v2.5.5", nil},
		{"a1111-meta.png", "a1111", "Fooocus v2.5.5", nil},
		{"a1111-meta.jpeg", "a1111", "Fooocus v2.5.5", nil},
		{"a1111-meta.webp", "a1111", "Fooocus v2.5.5", nil},
		{"2024-01-05_23-11-48_9167.png", "fooocus", "v2.1.860", metaV21Converted},
		{"2024-03-11_00-28-41_6907.png", "fooocus", "Fooocus v2.2.1", metaV22Converted},
	}

	assert.Len(t, images, len(testCases))

	for _, tc := range testCases {
		t.Run(tc.file, func(t *testing.T) {
			assert.Contains(t, images, tc.file, "Keys: %s", slices.Collect(maps.Keys(images)))

			metadata := images[tc.file]
			require.NotZero(t, metadata)
			require.IsType(t, Metadata{}, metadata)

			assert.Equal(t, tc.software, metadata.Version)
			assert.Equal(t, tc.scheme, metadata.MetadataScheme)
			// TODO reference data to compare to
			if tc.meta != nil {
				assert.Equal(t, *tc.meta, metadata)
			}
		})
	}
}
