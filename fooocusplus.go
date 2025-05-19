package metadata

import (
	"fmt"
	"io"
	"log/slog"
	"path"

	"github.com/fkleon/fooocus-metadata/internal/fooocusplus"
)

// FooocusPlusMetadataExtractor can decode embedded Fooocus Plus metadata from
// an image file, or external metadata from the Fooocus Plus private log.
type FooocusPlusMetadataExtractor struct {
	*baseMetadataExtractor
}

func (e FooocusPlusMetadataExtractor) Read(file ImageFile) (meta fooocusplus.Metadata, err error) {

	if file.exif != nil {
		return fooocusplus.ExtractMetadataFromExifData(file.exif)
	}

	if file.pngText != nil {
		return fooocusplus.ExtractMetadataFromPngData(file.pngText)
	}

	return meta, fmt.Errorf("missing EXIF or PNG tEXt")
}

func (e FooocusPlusMetadataExtractor) Extract(file ImageFile) (Parameters, error) {

	var param fooocusplus.Parameters

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

	if log, err := fooocusplus.ParsePrivateLog(logfile); err == nil {
		slog.Debug("Private log file", "file", logfile, "images", len(log))
		if meta, ok := log[file.Name()]; ok {
			param.Metadata = meta
			return param, nil
		}
	}

	return nil, fmt.Errorf("FooocusPlus: No metadata found")
}

func NewFooocusPlusMetadataExtractor() FooocusPlusMetadataExtractor {
	return FooocusPlusMetadataExtractor{
		baseMetadataExtractor: &baseMetadataExtractor{
			source:      FooocusPlus,
			dateLayout:  "2006-01-02_15-04-05",
			logfileName: "log.html",
		},
	}
}

func NewFooocusPlusMetadataReader() MetadataReader[fooocusplus.Metadata] {
	return NewFooocusPlusMetadataExtractor()
}

// FooocusPlusMetadataWriter can embed Fooocus Plus metadata into a PNG
// image file.
type FooocusPlusMetadataWriter struct {
	*baseMetadataWriter
}

func (w FooocusPlusMetadataWriter) Write(source io.Reader, target io.Writer, metadata fooocusplus.Metadata) error {
	values := map[string]interface{}{
		"Comment": metadata,
	}
	return w.embedMetadata(source, target, values)
}

func NewFooocusPlusMetadataWriter() MetadataWriter[fooocusplus.Metadata] {
	return FooocusPlusMetadataWriter{
		baseMetadataWriter: &baseMetadataWriter{
			target: FooocusPlus,
		},
	}
}

func init() {
	extractor := NewFooocusPlusMetadataExtractor()
	RegisterReader(extractor.source, extractor.Extract)
}
