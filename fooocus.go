package metadata

import (
	"fmt"
	"io"
	"log/slog"
	"path"

	"github.com/fkleon/fooocus-metadata/internal/fooocus"
)

// FooocusMetadataExtractor can decode embedded Fooocus metadata from
// an image file, or external metadata from the Fooocus private log.
type FooocusMetadataExtractor struct {
	*baseMetadataExtractor
}

func (e FooocusMetadataExtractor) Read(file ImageFile) (meta fooocus.Metadata, err error) {

	if file.exif != nil {
		return fooocus.ExtractMetadataFromExifData(file.exif)
	}

	if file.pngText != nil {
		return fooocus.ExtractMetadataFromPngData(file.pngText)
	}

	return meta, fmt.Errorf("missing EXIF or PNG tEXt")
}

func (e FooocusMetadataExtractor) Extract(file ImageFile) (Parameters, error) {

	var param fooocus.Parameters

	// Try to parse creation time from filename
	param.Created, _ = e.parseDateFromFilename(file.Name())

	slog.Debug("Checking embedded metadata..", "file", file.Name())
	if meta, err := e.Read(file); err == nil {
		param.Metadata = meta
		return param, nil
	}

	// Fallback to private log
	slog.Debug("Checking private log..", "logfile", e.logfileName)
	var logfile = path.Join(path.Dir(file.Path()), e.logfileName)

	if log, err := fooocus.ParsePrivateLog(logfile); err == nil {
		slog.Debug("Private log file", "file", logfile, "images", len(log))
		if meta, ok := log[file.Name()]; ok {
			param.Metadata = meta
			return param, nil
		}
	}

	return nil, fmt.Errorf("Fooocus: No metadata found")
}

func NewFooocusMetadataExtractor() FooocusMetadataExtractor {
	return FooocusMetadataExtractor{
		baseMetadataExtractor: &baseMetadataExtractor{
			source:      Fooocus,
			dateLayout:  "2006-01-02_15-04-05",
			logfileName: "log.html",
		},
	}
}

func NewFooocusMetadataReader() MetadataReader[fooocus.Metadata] {
	return NewFooocusMetadataExtractor()
}

// FooocusMetadataWriter can embed Fooocus metadata into a PNG
// image file.
type FooocusMetadataWriter struct {
	*baseMetadataWriter
}

func (w FooocusMetadataWriter) Write(source io.Reader, target io.Writer, metadata fooocus.Metadata) error {
	values := map[string]interface{}{
		"fooocus_scheme": fooocus.Fooocus.String(),
		"parameters":     metadata,
	}
	return w.embedMetadata(source, target, values)
}

func NewFooocusMetadataWriter() MetadataWriter[fooocus.Metadata] {
	return FooocusMetadataWriter{
		baseMetadataWriter: &baseMetadataWriter{
			target: Fooocus,
		},
	}
}

func init() {
	extractor := NewFooocusMetadataExtractor()
	RegisterReader(extractor.source, extractor.Extract)
}
