package metadata

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/fkleon/fooocus-metadata/internal/ruinedfooocus"
)

// RuinedFooocusMetadataExtractor can decode embedded Ruined Fooocus metadata from
// an image file.
type RuinedFooocusMetadataExtractor struct {
	*baseMetadataExtractor
}

func (e RuinedFooocusMetadataExtractor) Read(file ImageFile) (meta ruinedfooocus.Metadata, err error) {

	if file.pngText != nil {
		return ruinedfooocus.ExtractMetadataFromPngData(file.pngText)
	}

	return meta, fmt.Errorf("missing PNG tEXt")
}

func (e RuinedFooocusMetadataExtractor) Extract(file ImageFile) (Parameters, error) {

	var param ruinedfooocus.Parameters

	// Try to parse creation time from filename
	param.Created, _ = e.parseDateFromFilename(file.Name())

	slog.Debug("Checking embedded metadata..", "file", file.Name())
	if meta, err := e.Read(file); err == nil {
		param.Metadata = meta
		return param, nil
	}

	return nil, fmt.Errorf("RuinedFooocus: No metadata found")
}

func NewRuinedFooocusMetadataExtractor() RuinedFooocusMetadataExtractor {
	return RuinedFooocusMetadataExtractor{
		baseMetadataExtractor: &baseMetadataExtractor{
			source:     RuinedFooocus,
			dateLayout: "2006-01-02_15-04-05",
		},
	}
}

func NewRuinedFooocusMetadataReader() MetadataReader[ruinedfooocus.Metadata] {
	return NewRuinedFooocusMetadataExtractor()
}

// RuinedFooocusMetadataWriter can embed Ruined Fooocus metadata into a PNG
// image file.
type RuinedFooocusMetadataWriter struct {
	*baseMetadataWriter
}

func (w RuinedFooocusMetadataWriter) Write(source io.Reader, target io.Writer, metadata ruinedfooocus.Metadata) error {
	values := map[string]interface{}{
		"parameters": metadata,
	}
	return w.embedMetadata(source, target, values)
}

func NewRuinedFooocusMetadataWriter() MetadataWriter[ruinedfooocus.Metadata] {
	return RuinedFooocusMetadataWriter{
		baseMetadataWriter: &baseMetadataWriter{
			target: RuinedFooocus,
		},
	}
}

func init() {
	extractor := NewRuinedFooocusMetadataExtractor()
	RegisterReader(extractor.source, extractor.Extract)
}
