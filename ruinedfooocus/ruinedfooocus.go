package ruinedfooocus

import (
	"fmt"
	"io"
	"log/slog"
	"path/filepath"

	m "github.com/fkleon/fooocus-metadata/types"
)

// RuinedFooocusMetadataExtractor can decode embedded Ruined Fooocus metadata from
// an image file.
type RuinedFooocusMetadataExtractor struct {
	*m.FileMetadataExtractor
}

func (e RuinedFooocusMetadataExtractor) Decode(file m.ImageMetadataContext) (meta Metadata, err error) {

	var data = file.EmbeddedMetadata
	var parameters string

	// Parameters from PNG "parameters"
	if paramTag, ok := data["parameters"]; ok {
		parameters = paramTag.Value.(string)
	} else {
		return meta, fmt.Errorf("%s: Parameters not found", Software)
	}

	return parseMetadata(parameters)
}

func (e RuinedFooocusMetadataExtractor) Extract(file m.ImageMetadataContext) (m.StructuredMetadata, error) {

	var meta m.StructuredMetadata = m.StructuredMetadata{
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

func NewRuinedFooocusMetadataExtractor() m.Reader[Metadata] {
	return RuinedFooocusMetadataExtractor{
		FileMetadataExtractor: &m.FileMetadataExtractor{
			DateLayout: "2006-01-02_15-04-05",
		},
	}
}

// RuinedFooocusMetadataWriter can embed Ruined Fooocus metadata into a PNG
// image file.
type RuinedFooocusMetadataWriter struct {
	*m.PngMetadataWriter
}

func (w RuinedFooocusMetadataWriter) Write(target io.Writer, metadata Metadata) error {
	return w.CopyWrite(nil, target, metadata)
}

func (w RuinedFooocusMetadataWriter) CopyWrite(source io.Reader, target io.Writer, metadata Metadata) error {
	values := map[string]interface{}{
		"parameters": metadata,
	}
	return w.Embed(source, target, values)
}

func NewRuinedFooocusMetadataWriter() m.Writer[Metadata] {
	return RuinedFooocusMetadataWriter{
		PngMetadataWriter: &m.PngMetadataWriter{},
	}
}

func init() {
	extractor := NewRuinedFooocusMetadataExtractor()
	m.RegisterReader(Software, extractor.Extract)
}
