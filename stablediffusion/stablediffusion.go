package stablediffusion

import (
	"fmt"
	"log/slog"
	"path/filepath"

	m "github.com/fkleon/fooocus-metadata/types"
)

const (
	Software = "StableDiffusion"
)

// StableDiffusionMetadataExtractor can decode embedded A1111 metadata from
// an image file.
type StableDiffusionMetadataExtractor struct {
	*m.FileMetadataExtractor
}

func (e StableDiffusionMetadataExtractor) Decode(file m.ImageMetadataContext) (meta Metadata, err error) {

	var data = file.EmbeddedMetadata
	var parameters string

	// Parameters from PNG "parameters"
	if paramTag, ok := data["parameters"]; ok {
		parameters = paramTag.Value.(string)
	} else {
		return meta, fmt.Errorf("%s: Parameters not found", Software)
	}

	return ParseParameters(parameters)
}

func (e StableDiffusionMetadataExtractor) Extract(file m.ImageMetadataContext) (m.StructuredMetadata, error) {

	var meta = m.StructuredMetadata{
		Source: Software,
	}

	// Try to parse creation time from filename
	filename := filepath.Base(file.Filepath)
	meta.Created, _ = e.ParseDateFromFilename(filename)

	slog.Debug("Checking embedded metadata..", "file", filename)
	if params, err := e.Decode(file); err == nil {
		meta.Params = &Parameters{
			Metadata: params,
		}
		return meta, nil
	}

	return meta, fmt.Errorf("%s: No metadata found", Software)
}

func NewStableDiffusionMetadataExtractor() m.Reader[Metadata] {
	return StableDiffusionMetadataExtractor{
		FileMetadataExtractor: &m.FileMetadataExtractor{
			DateLayout: "2006-01-02_15-04-05",
		},
	}
}

func init() {
	extractor := NewStableDiffusionMetadataExtractor()
	m.RegisterReader(Software, extractor.Extract)
}
