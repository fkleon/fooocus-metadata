package fooocus

import (
	"time"

	"github.com/fkleon/fooocus-metadata/types"
)

// Adapter that implements the types.GenerationParameters
// interface on top of Fooocus Metadata.
type Parameters struct {
	Metadata
	Created time.Time
}

func (m Parameters) Version() string {
	return m.Metadata.Version
}

func (m Parameters) Model() string {
	return types.NormaliseModelName(m.BaseModel)
}

func (m Parameters) LoRAs() []types.Lora {
	var loras = make([]types.Lora, len(m.Loras))
	for i, lora := range m.Loras {
		loras[i] = types.Lora{
			Name:   types.NormaliseModelName(lora.Name),
			Weight: lora.Weight,
		}
	}
	return loras
}

func (m Parameters) PositivePrompt() string {
	return m.Prompt
}

func (m Parameters) NegativePrompt() string {
	return m.Metadata.NegativePrompt
}

func (m Parameters) Seed() string {
	return m.Metadata.Seed
}

func (m Parameters) CreatedTime() time.Time {
	return m.Created
}

func (m Parameters) Raw() interface{} {
	return m.Metadata
}
