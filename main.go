package metadata

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"

	"image"
	_ "image/jpeg"
	"image/png"

	_ "golang.org/x/image/webp"

	"github.com/fkleon/fooocus-metadata/internal/fooocus"
	"github.com/fkleon/fooocus-metadata/internal/fooocusplus"
	"github.com/fkleon/fooocus-metadata/internal/ruinedfooocus"
)

// Supported metadata sources.
//
//go:generate stringer -type=Software
type Software uint8

// Has returns true if the given software is set in the bitmap.
func (t Software) Has(software Software) bool {
	return t&software != 0
}

const (
	// Fooocus: https://github.com/lllyasviel/Fooocus
	Fooocus Software = 1 << iota
	// FooocusPlus: https://github.com/DavidDragonsage/FooocusPlus
	FooocusPlus
	// RuinedFoocus: https://github.com/runew0lf/RuinedFooocus
	RuinedFooocus
)

/*
	Process:

	- File
	- Get EXIF
	- Get tEXt
	- Get log?
	- Detect metadata type
	- From EXIF:
		- Software field?
		- version within parameters
	- From tEXt:
		- Software field?
		- schema field?
		- parameters or comment field?
		- version within parameters
	- From log:
		- version within parameters
	- Extract
*/

// ExtractOptions contains the options for the Extract function.
type ExtractOptions struct {
	// The path of the file to read image metadata from.
	File string

	// If set, will only consider the given software types.
	// Note that this is a bitmask and you may send multiple types at once.
	Sources Software
}

var extractors = [...]MetadataExtractor{
	NewFooocusMetadataExtractor(),
	NewFooocusPlusMetadataExtractor(),
	//NewRuinedFooocusMetadataExtractor(),
}

var r = NewFooocusMetadataReader()

type format struct {
	name   Software
	decode func(ImageFile) (Parameters, error)
}

type format2[M IMetadata] struct {
	name   Software
	decode func(ImageFile) (M, error)
	encode func(string, M) error
}

var formats []format = make([]format, 0, 3)

func RegisterReader(name Software, decode func(ImageFile) (Parameters, error)) {
	formats = append(formats, format{name, decode})
}

// TODO: ExtractFromFile returns..
func ExtractFromFile(filePath string) (params Parameters, err error) {
	return Extract(ExtractOptions{
		File:    filePath,
		Sources: Fooocus | FooocusPlus | RuinedFooocus,
	})
}

func Extract(opts ExtractOptions) (params Parameters, err error) {

	if opts.File == "" {
		return nil, fmt.Errorf("input file is required")
	}

	if opts.Sources == 0 {
		opts.Sources = Fooocus | FooocusPlus | RuinedFooocus
	}

	slog.Info("Extract", "options", opts)

	// Parse image metadata
	imageFile, err := OpenImageFile(opts.File)
	if err != nil {
		return
	}
	defer imageFile.Close()

	// These are the sources we support.
	sourceSet := Fooocus | FooocusPlus | RuinedFooocus
	// Remove sources that are not requested.
	sourceSet = sourceSet & opts.Sources

	for _, format := range formats {
		if sourceSet.Has(format.name) {
			slog.Info("Extracting metadata..", "software", format.name)
			if params, err = format.decode(*imageFile); err == nil {
				slog.Info("Found metadata", "software", params.Software())
				return
			}
		}
	}

	return nil, fmt.Errorf("no metadata found")
}

type IMetadata interface {
	fooocus.Metadata | fooocusplus.Metadata | ruinedfooocus.Metadata
}

func ExtractOne[M IMetadata](filepath string) (meta M, err error) {

	imageFile, err := OpenImageFile(filepath)
	if err != nil {
		return meta, err
	}
	defer imageFile.Close()

	switch m := any(meta).(type) {
	case fooocus.Metadata:
		reader := NewFooocusMetadataReader()
		mia, err := reader.Read(*imageFile)
		return any(mia).(M), err
	case fooocusplus.Metadata:
		reader := NewFooocusPlusMetadataReader()
		mia, err := reader.Read(*imageFile)
		return any(mia).(M), err
	case ruinedfooocus.Metadata:
		reader := NewRuinedFooocusMetadataReader()
		mia, err := reader.Read(*imageFile)
		return any(mia).(M), err
	default:
		return meta, fmt.Errorf("unsupported metadata type: %T", m)
	}
}

// EmbedOptions contains the options for the Embed function.
type EmbedOptions struct {
	// The path of the file to write image metadata to.
	Target io.Writer

	// The path of the file to copy data from.
	Source io.ReadSeeker
}

func EmbedIntoFile[M IMetadata](filepath string, meta M) error {

	target, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open target file for writing: %w", err)
	}
	defer target.Close()

	return Embed(EmbedOptions{
		Target: target,
	}, meta)
}

func convertToPng(in io.ReadSeeker) (out io.ReadSeeker, err error) {
	image, format, err := image.Decode(in)
	if err != nil {
		return
	}

	slog.Info("Decoded source image", "format", format)

	if format == "png" {
		in.Seek(0, 0)
		return in, nil
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, image)
	return bytes.NewReader(buf.Bytes()), err
}

func Embed[M IMetadata](opts EmbedOptions, meta M) (err error) {

	if opts.Target == nil {
		return fmt.Errorf("target file is required")
	}

	// If source is given, ensure source file is a valid PNG
	if opts.Source != nil {
		if opts.Source, err = convertToPng(opts.Source); err != nil {
			return fmt.Errorf("failed to read source file as PNG: %w", err)
		}
	} else {
		opts.Source = bytes.NewReader(pngTemplate)
	}

	slog.Info("Embed", "options", opts)

	switch m := any(meta).(type) {
	case fooocus.Metadata:
		writer := NewFooocusMetadataWriter()
		err = writer.Write(opts.Source, opts.Target, m)
	case fooocusplus.Metadata:
		writer := NewFooocusPlusMetadataWriter()
		err = writer.Write(opts.Source, opts.Target, m)
	case ruinedfooocus.Metadata:
		writer := NewRuinedFooocusMetadataWriter()
		err = writer.Write(opts.Source, opts.Target, m)
	}

	return
}
