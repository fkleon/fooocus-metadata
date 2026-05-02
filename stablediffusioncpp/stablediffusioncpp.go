package stablediffusioncpp

import (
	"fmt"
	"log/slog"
	"path/filepath"

	sd "github.com/fkleon/fooocus-metadata/stablediffusion"
	m "github.com/fkleon/fooocus-metadata/types"
)

const Software = "stable-diffusion.cpp"

// StableDiffusionCppMetadataExtractor can decode embedded SDCPP metadata from
// an image file.
type StableDiffusionCppMetadataExtractor struct {
	*m.FileMetadataExtractor
}

func (e StableDiffusionCppMetadataExtractor) Decode(file m.ImageMetadataContext) (meta sd.Metadata, err error) {
	data := file.EmbeddedMetadata
	var parameters string

	// Parameters from EXIF "UserComment" or PNG "parameters"
	if paramTag, ok := data["UserComment"]; ok {
		parameters = paramTag.Value.(string)
	} else if paramTag, ok := data["parameters"]; ok {
		parameters = paramTag.Value.(string)
	} else {
		return meta, fmt.Errorf("%s: Parameters not found", Software)
	}

	return ParseParameters(parameters)
}

func (e StableDiffusionCppMetadataExtractor) Extract(file m.ImageMetadataContext) (m.StructuredMetadata, error) {
	meta := m.StructuredMetadata{Source: Software}

	filename := filepath.Base(file.Filepath)
	meta.Created, _ = e.ParseDateFromFilename(filename)

	slog.Debug("Checking embedded metadata..", "file", filename)

	if params, err := e.Decode(file); err == nil {
		meta.Params = &sd.Parameters{Metadata: params}
		return meta, nil
	}

	return meta, fmt.Errorf("%s: No metadata found", Software)
}

func NewStableDiffusionCppMetadataExtractor() m.Reader[sd.Metadata] {
	return StableDiffusionCppMetadataExtractor{
		FileMetadataExtractor: &m.FileMetadataExtractor{
			DateLayouts: []string{
				"2006-01-02_15-04-05",
				"20060102_150405",
			},
		},
	}
}

func init() {
	extractor := NewStableDiffusionCppMetadataExtractor()
	m.RegisterReader(Software, extractor.Extract)
}
