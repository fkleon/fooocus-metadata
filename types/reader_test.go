package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDateNoLayout(t *testing.T) {
	reader := FileMetadataExtractor{}
	_, err := reader.ParseDateFromFilename("any.png")

	assert.Error(t, err)
}

func TestParseDate(t *testing.T) {
	reader := FileMetadataExtractor{
		DateLayouts: []string{
			"2006-01-02_15-04-05",
			"20060102_150405",
		},
	}

	expected := time.Date(2024, time.January, 5, 23, 11, 48, 0, time.UTC)
	testCases := []string{
		"2024-01-05_23-11-48.png",
		"2024-01-05_23-11-48_2.png",
		"20240105_231148.jpg",
		"20240105_231148_2.jpg",
		"20240105_231148_preview.jpg",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			actual, err := reader.ParseDateFromFilename(tc)
			require.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
	}
}
