package ruinedfooocus

import (
	"time"

	"github.com/fkleon/fooocus-metadata/types"
)

// Adapter that implements the types.GenerationParameters
// interface on top of RuinedFooocus Metadata.
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

func (m Parameters) LoRAs() []string {
	var loras []string = make([]string, len(m.Loras))
	for i, lora := range m.Loras {
		loras[i] = types.NormaliseModelName(lora.Name)
	}
	return loras
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
