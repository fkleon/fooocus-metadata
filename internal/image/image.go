// Package image provides utilities to detect image MIME types
// and read embedded image metadata (EXIF or PNG tEXt).
package image

import (
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/bep/imagemeta"
	pngembed "github.com/sabhiram/png-embed"
	"golang.org/x/text/encoding/charmap"

	"github.com/fkleon/fooocus-metadata/types"
)

// Wrapper around an os.File, adding MIME type detection
type File struct {
	*os.File
	MIME string
}

func (file *File) Path() string {
	return file.File.Name()
}
func (file *File) Name() string {
	return path.Base(file.Path())
}
func (file *File) Ext() string {
	return strings.ToLower(path.Ext(file.Path()))
}
func (file *File) IsImage() bool {
	return strings.HasPrefix(file.MIME, "image/")
}
func (file *File) detectMimeType() (err error) {
	// Rewind to the start
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	// Sniff content type via http.DetectContentType
	// Only the first 512 bytes are relevant.
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return
	}

	file.MIME = http.DetectContentType(buffer)

	// Fallback to file extension if MIME type could not be sniffed
	if file.MIME == "application/octet-stream" {
		mimeTypeByExt := mime.TypeByExtension(file.Ext())
		if mimeTypeByExt != "" {
			file.MIME = mimeTypeByExt
		}
	}

	slog.Debug("MIME type",
		"file", file.Name(),
		"mime", file.MIME,
	)

	return
}

func OpenFile(path string) (*File, error) {
	fin, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return NewFile(fin), nil
}

func NewFile(fin *os.File) *File {
	file := &File{
		File: fin,
	}
	_ = file.detectMimeType()
	return file
}

func NewContextFromReader(in io.ReadSeeker) (*types.ImageMetadataContext, error) {
	// Sniff MIME
	// Sniff content type via http.DetectContentType
	// Only the first 512 bytes are relevant.
	buffer := make([]byte, 512)
	_, err := in.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, err
	}

	mime := http.DetectContentType(buffer)
	return newContext(in, mime)
}

func NewContextFromFile(path string) (*types.ImageMetadataContext, error) {

	slog.Debug("Opening image file..", "filepath", path)

	file, err := OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if ctx, err := newContext(file, file.MIME); err == nil {
		ctx.Filepath = file.Path()
		return ctx, nil
	} else {
		return ctx, err
	}
}

func newContext(in io.ReadSeeker, mime string) (*types.ImageMetadataContext, error) {

	// Build image metadata and parse additional metadata sources
	var metadataMap map[string]imagemeta.TagInfo
	var metadataErr error

	switch mime {
	case "image/jpeg":
		fallthrough
	case "image/webp":
		fallthrough
	case "image/tiff":
		slog.Debug("Metadata source", "mime", mime, "source", "EXIF")
		var exif *imagemeta.Tags
		if exif, metadataErr = extractExif(in, mime); metadataErr == nil {
			metadataMap = exif.All()
		}
	case "image/png":
		slog.Debug("Metadata source", "mime", mime, "source", "PNG tEXt")
		var pngText map[string]string
		if pngText, metadataErr = extractPngText(in); metadataErr == nil {
			metadataMap = make(map[string]imagemeta.TagInfo, len(pngText))
			for k, v := range pngText {
				metadataMap[k] = imagemeta.TagInfo{
					Source:    0,
					Tag:       k,
					Namespace: "PNG/tEXt",
					Value:     v,
				}
			}
		}
	default:
		slog.Warn("Unsupported MIME type", "mime", mime)
		return nil, fmt.Errorf("unsupported MIME type: %s", mime)
	}

	if metadataErr != nil {
		slog.Warn("Failed to extract embedded metadata",
			"error", metadataErr)
	}

	return &types.ImageMetadataContext{
		MIME:             mime,
		EmbeddedMetadata: metadataMap,
	}, nil
}

func extractExif(fin io.ReadSeeker, mimeType string) (data *imagemeta.Tags, err error) {

	data = &imagemeta.Tags{}

	// Rewind to the start
	_, err = fin.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Extract EXIF data
	var format imagemeta.ImageFormat
	switch mimeType {
	default:
		fallthrough
	case "image/jpeg":
		format = imagemeta.JPEG
	case "image/webp":
		format = imagemeta.WebP
	case "image/tiff":
		format = imagemeta.TIFF
	}

	err = imagemeta.Decode(imagemeta.Options{
		R:           fin,
		ImageFormat: format,
		Sources:     imagemeta.EXIF,
		HandleTag: func(info imagemeta.TagInfo) error {
			data.Add(info)
			return nil
		},
		Warnf: func(msg string, args ...any) {
			slog.Debug(fmt.Sprintf("EXIF warning: %s", fmt.Sprintf(msg, args)))
		},
	})

	if err != nil {
		slog.Warn("Failed to extract EXIF data",
			"error", err)
	}

	return
}

func extractPngText(fin io.ReadSeeker) (pngText map[string]string, err error) {
	// Rewind to the start
	_, err = fin.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Extract PNG tEXt data
	pngText, err = extractPngTextChunks(fin)
	if err != nil {
		slog.Debug("Failed to extract PNG tEXt chunks",
			"error", err)
	}

	return
}

func extractPngTextChunks(fin io.ReadSeeker) (map[string]string, error) {

	data, err := io.ReadAll(fin)
	if err != nil {
		return nil, err
	}

	// Extract PNG tEXt chunks
	textData, err := pngembed.Extract(data)
	if err != nil {
		return nil, err
	}

	textDataDecoded := make(map[string]string)

	// Decode text with ISO-8859-1 as per PNG spec
	decoder := charmap.ISO8859_1.NewDecoder()
	for k, v := range textData {
		kd, err := decoder.String(string(k))
		if err != nil {
			slog.Warn("failed to decode key",
				"key", k,
				"error", err)
			continue
		}
		vd, err := decoder.String(string(v))
		if err != nil {
			slog.Warn("failed to decode value",
				"key", k, "value", v,
				"error", err)
			continue
		}
		textDataDecoded[kd] = vd
		textDataDecoded[k] = string(v)
	}

	return textDataDecoded, nil
}
