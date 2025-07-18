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
	DateLayout  string
	LogfileName string
}

func (e *FileMetadataExtractor) ParseDateFromFilename(filename string) (time.Time, error) {

	layoutIn := e.DateLayout

	if len(filename) < len(layoutIn) {
		return time.Time{}, fmt.Errorf("failed to parse date from filename: too short")
	}

	datepart := filename[:len(layoutIn)]
	return time.Parse(layoutIn, datepart)
}
