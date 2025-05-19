package metadata

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
	file.Seek(0, io.SeekStart)

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

// Wrapper around a File, adding image metadata extraction
type ImageFile struct {
	*File
	exif    *imagemeta.Tags
	pngText map[string]string
}

func OpenImageFile(path string) (image *ImageFile, err error) {

	slog.Debug("Opening image file..", "filepath", path)

	file, err := OpenFile(path)
	if err != nil {
		return
	}

	if !file.IsImage() {
		return nil, fmt.Errorf("file is not an image: %s", path)
	}

	// Build image metadata and parse additional metadata sources
	image = &ImageFile{
		File: file,
	}

	var metadataErr error

	switch image.MIME {
	case "image/jpeg":
		fallthrough
	case "image/webp":
		fallthrough
	case "image/tiff":
		slog.Debug("Metadata source", "mime", image.MIME, "source", "EXIF")
		image.exif, metadataErr = extractExif(image, image.MIME)
	case "image/png":
		slog.Debug("Metadata source", "mime", image.MIME, "source", "PNG tEXt")
		image.pngText, metadataErr = extractPngText(image)
	default:
		slog.Warn("Unsupported MIME type",
			"filepath", file.Path(), "mime", image.MIME)
		return
	}

	if metadataErr != nil {
		slog.Warn("Failed to extract embedded metadata",
			"filepath", file.Path(),
			"error", metadataErr)
	}

	return
}

func extractExif(fin io.ReadSeeker, mimeType string) (data *imagemeta.Tags, err error) {

	data = &imagemeta.Tags{}

	// Rewind to the start
	fin.Seek(0, io.SeekStart)

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

func extractPngText(fin io.ReadSeeker) (pngText map[string]string, pngErr error) {
	// Rewind to the start
	fin.Seek(0, io.SeekStart)

	// Extract PNG tEXt data
	pngText, pngErr = extractPngTextChunks(fin)
	if pngErr != nil {
		slog.Debug("Failed to extract PNG tEXt chunks",
			"error", pngErr)
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
