package fooocus

import (
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"

	_ "embed"

	"github.com/cozy/goexif2/exif"
	pngembed "github.com/sabhiram/png-embed"
	"golang.org/x/text/encoding/charmap"
)

//go:embed template.png
var pngTemplate []byte

type File struct {
	Path     string
	MIME     string
	FileInfo *os.FileInfo
	fin      *os.File
}

func (file *File) Name() string {
	return path.Base(file.Path)
}
func (file *File) Ext() string {
	return strings.ToLower(path.Ext(file.Path))
}
func (file *File) IsImage() bool {
	return strings.HasPrefix(file.MIME, "image/")
}
func (file *File) Parse(fin *os.File) (err error) {
	file.MIME, err = DetectMimeType(file, fin)
	stat, err := os.Stat(file.Path)
	if err != nil {
		file.FileInfo = &stat
	}
	return
}

type ImageFile struct {
	File
	FooocusMetadata *FooocusMeta
	exif            *exif.Exif
	pngText         map[string]string
}

func NewImageInfo(filePath string) (imageInfo *ImageFile, err error) {

	// Open and parse file metadata
	fin, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer fin.Close()

	file := &File{
		Path: filePath,
	}
	err = file.Parse(fin)
	if err != nil {
		return
	}

	if !file.IsImage() {
		err = fmt.Errorf("File is not an image: %s", filePath)
		return
	}

	// Build image metadata and parse additional metadata sources
	imageInfo = &ImageFile{
		File: *file,
	}

	var metadataErr error

	switch imageInfo.MIME {
	case "image/jpeg":
		fallthrough
	case "image/webp":
		fallthrough
	case "image/tiff":
		slog.Info("Extracting embedded metadata from EXIF..")
		if imageInfo.exif, metadataErr = ExtractExif(fin); metadataErr == nil {
			imageInfo.FooocusMetadata, metadataErr = ExtractMetadataFromExifData(imageInfo.exif)
		}
	case "image/png":
		slog.Info("Extracting embedded metadata from PNG tEXt..")
		if imageInfo.pngText, metadataErr = ExtractPngText(fin); metadataErr == nil {
			imageInfo.FooocusMetadata, metadataErr = ExtractMetadataFromPngData(imageInfo.pngText)
		}
	default:
		slog.Warn("Unsupported MIME type",
			"filepath", file.Path, "mime", imageInfo.MIME)
		return
	}

	if metadataErr != nil {
		slog.Warn("Failed to extract embedded metadata",
			"filepath", file.Path,
			"error", metadataErr)
	}

	return
}

func DetectMimeType(file *File, fin *os.File) (mimeType string, err error) {
	// Rewind to the start
	fin.Seek(0, io.SeekStart)

	// Sniff content type via http.DetectContentType
	// Only the first 512 bytes are relevant.
	buffer := make([]byte, 512)
	_, err = fin.Read(buffer)
	if err != nil && err != io.EOF {
		return
	}

	mimeType = http.DetectContentType(buffer)

	// Fallback to file extension if MIME type could not be sniffed
	if mimeType == "application/octet-stream" {
		mimeType = mime.TypeByExtension(file.Ext())
	}

	slog.Debug("Detected MIME type",
		"file", file.Name(),
		"mime", mimeType,
	)

	return
}

func ExtractExif(fin *os.File) (exifData *exif.Exif, exifErr error) {
	// Rewind to the start
	fin.Seek(0, io.SeekStart)

	// Extract EXIF data
	exifData, exifErr = exif.Decode(fin)
	if exifErr != nil {
		slog.Debug("Failed to extract EXIF data",
			"error", exifErr)
	}

	return
}

func ExtractPngText(fin *os.File) (pngText map[string]string, pngErr error) {
	// Rewind to the start
	fin.Seek(0, io.SeekStart)

	// Extract PNG tEXt data
	pngText, pngErr = extractPngTextChunks(fin)
	if pngErr != nil {
		slog.Debug("Failed to extract PNG tEXT chunks",
			"error", pngErr)
	}

	return
}

func extractPngTextChunks(fin *os.File) (map[string]string, error) {

	data, err := os.ReadFile(fin.Name())
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

func EmbedMetadataAsPngText(source *os.File, target *os.File, meta *FooocusMeta) (err error) {

	slog.Debug("Embedding Fooocus metadata into PNG",
		"target", target.Name())

	var data []byte

	if source != nil {
		data, err = os.ReadFile(source.Name())
		if err != nil {
			return err
		}
	} else {
		data = pngTemplate
	}

	kvs := map[string]interface{}{
		"fooocus_scheme": fooocus,
		"parameters":     meta,
	}

	for k, v := range kvs {
		data, err = pngembed.Embed(data, k, v)
		if err != nil {
			return err
		}
	}

	_, err = target.Write(data)
	return err
}
