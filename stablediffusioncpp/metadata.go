// Package stablediffusioncpp implements reading [stable-diffusion.cpp] metadata
// (image generation parameters).
//
// [stable-diffusion.cpp]: https://github.com/leejet/stable-diffusion.cpp
package stablediffusioncpp

import (
	"encoding/json"
	"fmt"
	"strings"

	sd "github.com/fkleon/fooocus-metadata/stablediffusion"
)

type metadata struct {
	Schema     string  `json:"schema"`
	Seed       int     `json:"seed"`
	Width      int     `json:"width"`
	Height     int     `json:"height"`
	ClipSkip   int     `json:"clip_skip"`
	Rng        string  `json:"rng"`
	SamplerRng string  `json:"sampler_rng"`
	Strength   float32 `json:"strength"`
	Generator  struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Commit  string `json:"commit"`
	} `json:"generator"`
	Prompt struct {
		Positive string `json:"positive"`
		Negative string `json:"negative"`
	} `json:"prompt"`
	Sampling struct {
		Steps        int       `json:"steps"`
		Eta          float32   `json:"eta"`
		Method       string    `json:"method"`
		Scheduler    string    `json:"scheduler"`
		CustomSigmas []float64 `json:"custom_sigmas"`
		Guidance     struct {
			TxtCfg            float32 `json:"txt_cfg"`
			DistilledGuidance float32 `json:"distilled_guidance"`
			SLG               struct {
				Scale  float32 `json:"scale"`
				Layers []int   `json:"layers"`
				Start  float32 `json:"start"`
				End    float32 `json:"end"`
			} `json:"slg"`
		} `json:"guidance"`
	} `json:"sampling"`
	Models struct {
		Model          string `json:"model"`
		ClipL          string `json:"clip_l"`
		ClipG          string `json:"clip_g"`
		T5XXL          string `json:"t5xxl"`
		LLM            string `json:"llm"`
		LLMVision      string `json:"llm_vision"`
		DiffusionModel string `json:"diffusion_model"`
		VAE            string `json:"vae"`
	} `json:"models"`
	Loras []struct {
		Name       string  `json:"name"`
		Multiplier float32 `json:"multiplier"`
	} `json:"loras"`
	Hires *struct {
		Upscaler          string  `json:"upscaler"`
		Scale             float32 `json:"scale"`
		TargetWidth       int     `json:"target_width"`
		TargetHeight      int     `json:"target_height"`
		Steps             int     `json:"steps"`
		DenoisingStrength float32 `json:"denoising_strength"`
	} `json:"hires"`
}

func ParseParameters(in string) (sd.Metadata, error) {
	payload, err := extractPayload(in)
	if err != nil {
		return sd.Metadata{}, err
	}

	var raw metadata
	if err := json.Unmarshal(payload, &raw); err != nil {
		return sd.Metadata{}, fmt.Errorf("invalid SDCPP metadata: %w", err)
	}

	if !strings.HasPrefix(raw.Schema, "sdcpp.image.params/") {
		return sd.Metadata{}, fmt.Errorf("unsupported SDCPP schema %q", raw.Schema)
	}

	meta := sd.Metadata{
		CfgScale:          raw.Sampling.Guidance.TxtCfg,
		ClipSkip:          raw.ClipSkip,
		CustomSigmas:      toFloatSlice(raw.Sampling.CustomSigmas),
		DenoisingStrength: raw.Strength,
		Guidance:          raw.Sampling.Guidance.DistilledGuidance,
		Model:             raw.Models.Model,
		NegativePrompt:    raw.Prompt.Negative,
		Prompt:            raw.Prompt.Positive,
		Rng:               raw.Rng,
		SamplerRng:        raw.SamplerRng,
		Sampler:           samplerName(raw.Sampling.Method, raw.Sampling.Scheduler, len(raw.Sampling.CustomSigmas) > 0),
		Seed:              raw.Seed,
		Steps:             raw.Sampling.Steps,
		Unet:              raw.Models.DiffusionModel,
		Vae:               raw.Models.VAE,
		Version:           version(raw.Generator.Name, raw.Generator.Version, raw.Generator.Commit),
	}

	if raw.Sampling.Guidance.SLG.Scale != 0 || len(raw.Sampling.Guidance.SLG.Layers) > 0 {
		meta.SkipLayers = raw.Sampling.Guidance.SLG.Layers
		meta.SkipLayerEnd = raw.Sampling.Guidance.SLG.End
		meta.SkipLayerStart = raw.Sampling.Guidance.SLG.Start
		meta.SLGScale = raw.Sampling.Guidance.SLG.Scale
	}

	if raw.Sampling.Eta != 0.0 {
		meta.Eta = sd.Float(raw.Sampling.Eta)
	}

	if raw.Width > 0 || raw.Height > 0 {
		meta.Size = &sd.Size{Width: raw.Width, Height: raw.Height}
	}

	if te := strings.Join(compactStrings(
		raw.Models.ClipL,
		raw.Models.ClipG,
		raw.Models.T5XXL,
		raw.Models.LLM,
		raw.Models.LLMVision,
	), ", "); te != "" {
		meta.TextEncoder = te
	}

	if len(raw.Loras) > 0 {
		meta.Loras = make(sd.Loras, 0, len(raw.Loras))
		for _, lora := range raw.Loras {
			meta.Loras = append(meta.Loras, sd.Lora{Name: lora.Name, Weight: lora.Multiplier})
		}
	}

	if raw.Hires != nil {
		meta.HiresUpscaler = raw.Hires.Upscaler
		meta.HiresUpscale = raw.Hires.Scale
		meta.HiresSteps = raw.Hires.Steps

		meta.DenoisingStrength = raw.Hires.DenoisingStrength
		if raw.Hires.TargetWidth > 0 || raw.Hires.TargetHeight > 0 {
			meta.HiresResize = &sd.Size{Width: raw.Hires.TargetWidth, Height: raw.Hires.TargetHeight}
		}
	}

	return meta, nil
}

func extractPayload(in string) ([]byte, error) {
	trimmed := strings.TrimSpace(in)
	if json.Valid([]byte(trimmed)) {
		return []byte(trimmed), nil
	}

	const field = "SDCPP:"

	idx := strings.LastIndex(in, field)
	if idx == -1 {
		return nil, fmt.Errorf("SDCPP field not found")
	}

	value := strings.TrimSpace(in[idx+len(field):])
	if value == "" {
		return nil, fmt.Errorf("SDCPP field is empty")
	}

	decoder := json.NewDecoder(strings.NewReader(value))

	var payload json.RawMessage
	if err := decoder.Decode(&payload); err != nil {
		return nil, fmt.Errorf("invalid SDCPP JSON: %w", err)
	}

	if rest := strings.TrimSpace(value[int(decoder.InputOffset()):]); rest != "" {
		return nil, fmt.Errorf("unexpected trailing data after SDCPP JSON")
	}

	return payload, nil
}

func samplerName(method string, scheduler string, hasCustomSigmas bool) string {
	parts := compactStrings(method)
	if !hasCustomSigmas {
		parts = append(parts, compactStrings(scheduler)...)
	}

	return strings.Join(parts, " ")
}

func version(name string, value string, commit string) string {
	switch {
	case value != "" && value != "unknown":
		return fmt.Sprintf("%s (%s)", name, value)
	case commit != "":
		return fmt.Sprintf("%s (%s)", name, commit)
	default:
		return name
	}
}

func compactStrings(values ...string) []string {
	compact := make([]string, 0, len(values))
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			compact = append(compact, value)
		}
	}

	return compact
}

func toFloatSlice(values []float64) sd.FloatSlice {
	if len(values) == 0 {
		return nil
	}

	converted := make(sd.FloatSlice, 0, len(values))
	for _, value := range values {
		converted = append(converted, float32(value))
	}

	return converted
}
