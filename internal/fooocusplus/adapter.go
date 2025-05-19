package fooocusplus

import "time"

// Adapter
type Parameters struct {
	Metadata
	Created time.Time
}

func (m Parameters) Software() string {
	return m.Metadata.Version
}

func (m Parameters) Model() string {
	return m.Metadata.BaseModel
}

func (m Parameters) Prompt() string {
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
