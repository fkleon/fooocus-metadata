package types

import (
	"path/filepath"
	"strings"
	"time"
)

type StructuredMetadata struct {
	// Source is the identifier for the tool that generated the metadata,
	// e.g., "Fooocus", "FooocusPlus" or "RuinedFooocus".
	Source string

	Created time.Time
	Params  GenerationParameters
}

type GenerationParameters interface {

	// The version of the software that generated the metadata.
	Version() string

	// The prompt used for the generation.
	PositivePrompt() string
	// The negative prompt used for the generation.
	NegativePrompt() string
	// The model used for the generation.
	Model() string
	// The LoRAs used for the generation.
	LoRAs() []string

	// Raw returns the underlying metadata struct (e.g. fooocus.Metadata).
	// The caller can type-assert it if needed.
	Raw() interface{}
}

func NormaliseModelName(name string) string {
	// Normalise model name by removing path
	// and known file extension
	base := filepath.Base(name)
	ext := filepath.Ext(name)

	switch strings.ToLower(ext) {
	case ".safetensors", ".gguf", ".pt", ".pth", ".onnx":
		return strings.TrimSuffix(base, ext)
	default:
		// Unknown extension, return as is
		return base
	}
}
