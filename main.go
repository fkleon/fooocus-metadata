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
// parsing creation date from file pattern or sidecase support
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
