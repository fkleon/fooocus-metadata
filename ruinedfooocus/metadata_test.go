package ruinedfooocus

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Metadata as JSON
const metaJson = `{
  "Prompt": "cinematic film still A sunflower field, shallow depth of field, vignette, highly detailed, high budget Hollywood film, cinemascope, moody, epic, gorgeous, ",
  "Negative": "anime, cartoon, graphic, text, painting, crayon, graphite, abstract, glitch, blur, bokeh , , ",
  "steps": 30,
  "cfg": 8.5,
  "width": 1152,
  "height": 896,
  "seed": 3864674281,
  "sampler_name": "dpmpp_2m_sde_gpu",
  "scheduler": "karras",
  "base_model_name": "sd_xl_base_1.0_0.9vae.safetensors",
  "base_model_hash": "be9edd61",
  "loras": [
    [
      "4852686128",
      "0.1 - sd_xl_offset_example-lora_1.0.safetensors"
    ]
  ],
  "start_step": 0,
  "denoise": null,
  "clip_skip": 1,
  "software": "RuinedFooocus"
}`

var meta = &Metadata{
	BaseModel:     "sd_xl_base_1.0_0.9vae.safetensors",
	BaseModelHash: "be9edd61",
	CfgScale:      8.5,
	ClipSkip:      1,
	Denoise:       nil,
	Height:        896,
	Width:         1152,
	Loras: []Lora{{
		Hash:   "4852686128",
		Name:   "sd_xl_offset_example-lora_1.0.safetensors",
		Weight: 0.1,
	}},
	NegativePrompt: "anime, cartoon, graphic, text, painting, crayon, graphite, abstract, glitch, blur, bokeh , , ",
	Prompt:         "cinematic film still A sunflower field, shallow depth of field, vignette, highly detailed, high budget Hollywood film, cinemascope, moody, epic, gorgeous, ",
	Sampler:        "dpmpp_2m_sde_gpu",
	Scheduler:      "karras",
	Seed:           3864674281,
	StartStep:      0,
	Steps:          30,
	Version:        "RuinedFooocus",
}

func TestDecodeMetadata(t *testing.T) {
	var decoded *Metadata
	err := json.Unmarshal([]byte(metaJson), &decoded)
	require.NoError(t, err)
	assert.Equal(t, meta, decoded)
}

func TestEncodeMetadata(t *testing.T) {
	encoded, err := json.Marshal(meta)
	require.NoError(t, err)
	assert.JSONEq(t, metaJson, string(encoded))
}
