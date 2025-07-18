package fooocusplus

import (
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"strings"

	m "github.com/fkleon/fooocus-metadata/types"
)

// FooocusPlusMetadataExtractor can decode embedded Fooocus Plus metadata from
// an image file, or external metadata from the Fooocus Plus private log.
type FooocusPlusMetadataExtractor struct {
	*m.FileMetadataExtractor
}

func (e FooocusPlusMetadataExtractor) Decode(file m.ImageMetadataContext) (meta Metadata, err error) {

	// TODO: scheme 'simple' if 'Comment' field exists
	var data = file.EmbeddedMetadata
	var parameters string

	// Software version from EXIF "Software"
	if softwareTag, ok := data["Software"]; ok {
		softwareVersion := softwareTag.Value.(string)
		if !strings.HasPrefix(softwareVersion, "FooocusPlus 1.") {
			return meta, fmt.Errorf("%s: EXIF: Unsupported software: %s", Software, softwareVersion)
		}
	}

	// Parameters from EXIF "UserComment" or PNG "Comment"
	if paramTag, ok := data["UserComment"]; ok {
		parameters = paramTag.Value.(string)
	} else {
		if paramTag, ok := data["Comment"]; ok {
			parameters = paramTag.Value.(string)
		} else {
			return meta, fmt.Errorf("%s: Parameters not found", Software)
		}
	}

	return parseMetadata(parameters)
}

func (e FooocusPlusMetadataExtractor) Extract(file m.ImageMetadataContext) (m.StructuredMetadata, error) {

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

	// Fallback to private log
	slog.Debug("Checking private log..", "logfile", e.LogfileName)
	var logfile = filepath.Join(filepath.Dir(file.Filepath), e.LogfileName)

	if log, err := ParsePrivateLog(logfile); err == nil {
		slog.Debug("Private log file", "file", logfile, "images", len(log))
		if params, ok := log[filename]; ok {
			meta.Params = &Parameters{
				Metadata: params,
			}
			return meta, nil
		}
	}

	return meta, fmt.Errorf("%s: No metadata found", Software)
}

func NewFooocusPlusMetadataExtractor() m.Reader[Metadata] {
	return FooocusPlusMetadataExtractor{
		FileMetadataExtractor: &m.FileMetadataExtractor{
			DateLayout:  "2006-01-02_15-04-05",
			LogfileName: "log.html",
		},
	}
}

// FooocusPlusMetadataWriter can embed Fooocus Plus metadata into a PNG
// image file.
type FooocusPlusMetadataWriter struct {
	*m.PngMetadataWriter
}

func (w FooocusPlusMetadataWriter) Write(target io.Writer, metadata Metadata) error {
	return w.CopyWrite(nil, target, metadata)
}

func (w FooocusPlusMetadataWriter) CopyWrite(source io.Reader, target io.Writer, metadata Metadata) error {
	values := map[string]interface{}{
		"Comment": metadata,
	}
	return w.Embed(source, target, values)
}

func NewFooocusPlusMetadataWriter() m.Writer[Metadata] {
	return FooocusPlusMetadataWriter{
		PngMetadataWriter: m.NewPngMetadataWriter(),
	}
}

func init() {
	extractor := NewFooocusPlusMetadataExtractor()
	m.RegisterReader(Software, extractor.Extract)
}
