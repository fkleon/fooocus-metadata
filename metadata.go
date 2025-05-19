package metadata

import (
	"fmt"
	"io"
	"log/slog"
	"time"

	_ "embed"

	pngembed "github.com/sabhiram/png-embed"
)

type Metadata struct {
	Software Software
	Created  time.Time

	// One of [fooocus.Metadata], [fooocusplus.Metadata] or [ruinedfooocus.Metadata].
	Raw interface{}
}

type Parameters interface {
	Software() string

	Model() string
	Prompt() string
	NegativePrompt() string

	CreatedTime() time.Time

	Raw() interface{}
}

type MetadataExtractor interface {
	Source() Software
	Extract(file ImageFile) (Parameters, error)
}

// MetadataReader reads embedded metadata of type M from a file.
type MetadataReader[M IMetadata] interface {
	Source() Software
	Read(file ImageFile) (M, error)
}

// MetadataWriter writes metadata of type M to a file.
type MetadataWriter[M IMetadata] interface {
	Target() Software
	Write(source io.Reader, target io.Writer, metadata M) error
}

// baseMetadataExtractor is a common base for metadata extractors.
type baseMetadataExtractor struct {
	source      Software
	dateLayout  string
	logfileName string
}

func (e *baseMetadataExtractor) Source() Software {
	return e.source
}

func (e *baseMetadataExtractor) parseDateFromFilename(filename string) (time.Time, error) {

	layoutIn := e.dateLayout

	if len(filename) < len(layoutIn) {
		return time.Time{}, fmt.Errorf("failed to parse date from filename: too short")
	}

	datepart := filename[:len(layoutIn)]
	return time.Parse(layoutIn, datepart)
}

//go:embed template.png
var pngTemplate []byte

// baseMetadataWriter is a common base for metadata writers.
type baseMetadataWriter struct {
	target Software
}

func (e *baseMetadataWriter) Target() Software {
	return e.target
}

func (e *baseMetadataWriter) embedMetadata(source io.Reader, target io.Writer, values map[string]interface{}) (err error) {

	slog.Debug("Embedding metadata", "count", len(values), "target", target)

	if source == nil {
		return fmt.Errorf("source is required")
	}

	if target == nil {
		return fmt.Errorf("target is required")
	}

	data, err := io.ReadAll(source)
	if err != nil {
		return fmt.Errorf("failed to read source: %w", err)
	}

	for k, v := range values {
		data, err = pngembed.Embed(data, k, v)
		if err != nil {
			return err
		}
	}

	_, err = target.Write(data)
	return err
}
