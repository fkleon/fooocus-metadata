package fooocusplus

import (
	"path"
	"strings"
	"time"
)

// Adapter that implements the types.GenerationParameters
// interface on top of FooocusPlus Metadata.
type Parameters struct {
	Metadata
	Created time.Time
}

func (m Parameters) Version() string {
	return m.Metadata.Version
}

func (m Parameters) Model() string {
	// Normalise model name by removing file extension
	return strings.TrimSuffix(m.BaseModel, path.Ext(m.BaseModel))
}

func (m Parameters) PositivePrompt() string {
	return m.Metadata.Prompt
}

func (m Parameters) NegativePrompt() string {
	return m.Metadata.NegativePrompt
}

func (m Parameters) CreatedTime() time.Time {
	return m.Created
}

func (m Parameters) Raw() interface{} {
	return m.Metadata
}
