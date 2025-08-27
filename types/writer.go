package types

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/png"
	"io"
	"log/slog"

	pngembed "github.com/sabhiram/png-embed"
)

//go:embed template.png
var PngTemplate []byte

// Writer is a generic interface for writing metadata.
type Writer[M any] interface {
	Write(target io.Writer, metadata M) error
	CopyWrite(source io.Reader, target io.Writer, metadata M) error
}

// PngMetadataWriter is a common base for metadata writers
// that embed into PNG.
type PngMetadataWriter struct {
	Template io.Reader
}

func NewPngMetadataWriter() *PngMetadataWriter {
	return &PngMetadataWriter{
		Template: bytes.NewReader(PngTemplate),
	}
}

func (e *PngMetadataWriter) Embed(source io.Reader, target io.Writer, values map[string]interface{}) (err error) {

	slog.Debug("Embedding metadata", "count", len(values), "target", target)

	if source == nil && e.Template == nil {
		return fmt.Errorf("source or template are required")
	}

	if target == nil {
		return fmt.Errorf("target is required")
	}

	if source == nil {
		source = e.Template
	} else {
		if source, err = convertToPng(source); err != nil {
			return fmt.Errorf("failed to convert source to PNG: %w", err)
		}
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

func convertToPng(in io.Reader) (out io.ReadSeeker, err error) {
	image, format, err := image.Decode(in)
	if err != nil {
		return
	}

	slog.Debug("Decoded source image", "format", format)

	buf := new(bytes.Buffer)
	err = png.Encode(buf, image)
	return bytes.NewReader(buf.Bytes()), err
}
