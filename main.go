// Package metadata provides high-level functions to extract
// image generation parameters from image files.
//
// This supports multiple sources, including Fooocus, FooocusPlus
// and RuinedFooocus.
//
// Example usage:
//
// To read metadata, first register the sources you want to enable
// by importing the appropriate package:
//
//	import (
//	  _ "github.com/fkleon/fooocus-metadata/fooocus"
//	)
//
// To read metadata from a file, use ExtractFromFile:
//
//	path := "fooocus/testdata/fooocus-meta.png"
//	meta, err := ExtractFromFile(path)
//	fmt.Println(meta.Version) // prints "Fooocus v2.5.5"
//
// To read metadata from a stream, use ExtractFromReader.
// You can provide a hint about the original filepath to
// enable path-based features such as sidecar support:
//
//	path := "fooocus/testdata/fooocus-meta.png"
//	file, err := os.Open(path)
//	defer file.Close()
//	meta, err := ExtractFromReader(file, WithPath(path))
//	fmt.Println(meta.Version) // prints "Fooocus v2.5.5"
//
// To write metadata, use the individual metadata writers
// provided by each package.
package metadata

import (
	"fmt"
	"io"
	"log/slog"

	// Required image decoders
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"

	"github.com/fkleon/fooocus-metadata/internal/image"
	"github.com/fkleon/fooocus-metadata/types"
)

// ExtractOptions contains the options for the Extract function.
type ExtractOptions struct {
	// The path of the file to read image metadata from.
	File string
}

type Config struct {
	Path string
}
type Option func(*Config)

// To enable path-based extraction features, e.g.
// parsing creation date from file pattern or sidecar support
func WithPath(path string) Option {
	return func(cfg *Config) {
		cfg.Path = path
	}
}

func ExtractFromFile(path string) (params types.StructuredMetadata, err error) {
	slog.Info("ExtractFromFile", "path", path)

	imageFile, err := image.NewContextFromFile(path)
	if err != nil {
		return
	}

	return types.Decode(*imageFile)
}

func ExtractFromReader(reader io.ReadSeeker, opts ...Option) (params types.StructuredMetadata, err error) {
	if reader == nil {
		return params, fmt.Errorf("input reader is required")
	}

	// Customise config
	cfg := Config{}
	for _, opt := range opts {
		opt(&cfg)
	}

	slog.Info("ExtractFromReader", "options", opts)

	// Parse image metadata
	imageCtx, err := image.NewContextFromReader(reader)
	if err != nil {
		return
	}
	imageCtx.Filepath = cfg.Path

	return types.Decode(*imageCtx)
}
