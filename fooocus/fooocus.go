package fooocus

import (
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"strings"

	m "github.com/fkleon/fooocus-metadata/types"
)

// FooocusMetadataExtractor can decode embedded Fooocus metadata from
// an image file, or external metadata from the Fooocus private log.
type FooocusMetadataExtractor struct {
	*m.FileMetadataExtractor
}

func (e FooocusMetadataExtractor) Decode(file m.ImageMetadataContext) (meta Metadata, err error) {

	var data = file.EmbeddedMetadata
	var scheme, parameters string

	// Software version from EXIF "Software"
	if softwareTag, ok := data["Software"]; ok {
		softwareVersion := softwareTag.Value.(string)
		if !strings.HasPrefix(softwareVersion, "Fooocus ") {
			return meta, fmt.Errorf("%s: EXIF: Unsupported software: %s", Software, softwareVersion)
		}
	}

	// Schema from EXIF "MakerNoteApple" or PNG "fooocus_scheme"
	if schemeTag, ok := data["MakerNoteApple"]; ok {
		scheme = schemeTag.Value.(string)
	} else {
		if schemeTag, ok := data["fooocus_scheme"]; ok {
			scheme = schemeTag.Value.(string)
		} else {
			return meta, fmt.Errorf("%s: Scheme not found", Software)
		}
	}

	// Parameters from EXIF "UserComment" or PNG "parameters"
	if paramTag, ok := data["UserComment"]; ok {
		parameters = paramTag.Value.(string)
	} else {
		if paramTag, ok := data["parameters"]; ok {
			parameters = paramTag.Value.(string)
		} else {
			return meta, fmt.Errorf("%s: Parameters not found", Software)
		}
	}

	return parseMetadata(scheme, parameters)
}

func (e FooocusMetadataExtractor) Extract(file m.ImageMetadataContext) (m.StructuredMetadata, error) {

	var meta = m.StructuredMetadata{
		Source: Software,
	}

	// Try to parse creation time from filename
	filename := filepath.Base(file.Filepath)
	meta.Created, _ = e.ParseDateFromFilename(filename)

	slog.Debug("Checking embedded metadata..", "file", file.Filepath)
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

func NewFooocusMetadataExtractor() m.Reader[Metadata] {
	return FooocusMetadataExtractor{
		FileMetadataExtractor: &m.FileMetadataExtractor{
			DateLayout:  "2006-01-02_15-04-05",
			LogfileName: "log.html",
		},
	}
}

// FooocusMetadataWriter can embed Fooocus metadata into a PNG
// image file.
type FooocusMetadataWriter struct {
	*m.PngMetadataWriter
}

func (w FooocusMetadataWriter) Write(target io.Writer, metadata Metadata) error {
	return w.CopyWrite(nil, target, metadata)
}

func (w FooocusMetadataWriter) CopyWrite(source io.Reader, target io.Writer, metadata Metadata) error {
	values := map[string]interface{}{
		"fooocus_scheme": Fooocus.String(),
		"parameters":     metadata,
	}
	return w.Embed(source, target, values)
}

func NewFooocusMetadataWriter() m.Writer[Metadata] {
	return FooocusMetadataWriter{
		PngMetadataWriter: m.NewPngMetadataWriter(),
	}
}

func init() {
	extractor := NewFooocusMetadataExtractor()
	m.RegisterReader(Software, extractor.Extract)
}
