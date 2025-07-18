package types

import "time"

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

	// Raw returns the underlying metadata struct (e.g. fooocus.Metadata).
	// The caller can type-assert it if needed.
	Raw() interface{}
}
