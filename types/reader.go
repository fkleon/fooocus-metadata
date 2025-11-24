package types

import (
	"fmt"
	"time"
)

// Reader is an generic interface for extracting metadata.
type Reader[T any] interface {
	// Decode reads software-specific metadata from the image.
	Decode(ImageMetadataContext) (T, error)
	// Extract reads structured metadata for the image.
	// This could come from embedded metadata or an external source
	// such as a sidecar file.
	Extract(ImageMetadataContext) (StructuredMetadata, error)
}

// FileMetadataExtractor is a common base for file-based metadata extractors.
type FileMetadataExtractor struct {
	DateLayouts []string
	LogfileName string
}

func (e *FileMetadataExtractor) ParseDateFromFilename(filename string) (time.Time, error) {

	for _, layoutIn := range e.DateLayouts {
		if len(filename) < len(layoutIn) {
			continue
		}

		datepart := filename[:len(layoutIn)]
		t, err := time.Parse(layoutIn, datepart)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("failed to parse date from filename: %s", filename)
}
