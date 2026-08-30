// Package stablediffusion implements reading [AUTOMATIC1111] style plaintext metadata
// (image generation parameters).
//
// [AUTOMATIC1111]: https://github.com/AUTOMATIC1111/stable-diffusion-webui
package stablediffusion

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type Metadata struct {
	BatchSize            int        `json:"batch_size,string,omitempty"`
	BatchPos             int        `json:"batch_pos,string,omitempty"`
	CfgScale             float32    `json:"cfg_scale,string"`
	ClipSkip             int        `json:"clip_skip,string,omitempty"`
	CustomSigmas         FloatSlice `json:"custom_sigmas,omitempty"`
	DenoisingStrength    float32    `json:"denoising_strength,string,omitempty"`
	Eta                  Float      `json:"eta,omitempty"`
	HiresSteps           int        `json:"hires_steps,string,omitempty"`
	HiresUpscale         float32    `json:"hires_upscale,string,omitempty"`
	HiresUpscaler        string     `json:"hires_upscaler,omitempty"`
	Guidance             float32    `json:"guidance,string,omitempty"`
	ImageNoiseMultiplier float32    `json:"image_noise_multiplier,string,omitempty"`
	Loras                Loras      `json:"loras,omitempty"`
	Model                string     `json:"model,omitempty"`
	ModelHash            string     `json:"model_hash,omitempty"`
	NegativePrompt       string     `json:"negative_prompt,omitempty"`
	Prompt               string     `json:"prompt"`
	Rng                  string     `json:"rng,omitempty"`
	SamplerRng           string     `json:"sampler_rng,omitempty"`
	Sampler              string     `json:"sampler"` // The Sampler field contains both the sampler and scheduler names
	Seed                 int        `json:"seed,string"`
	SkipLayers           IntSlice   `json:"skip_layers,omitempty"`
	SkipLayerEnd         float32    `json:"skip_layer_end,string,omitempty"`
	SkipLayerStart       float32    `json:"skip_layer_start,string,omitempty"`
	SLGScale             float32    `json:"slg_scale,string,omitempty"`
	Size                 *Size      `json:"size,omitempty"`
	Steps                int        `json:"steps,string"`
	TextEncoder          string     `json:"TE,omitempty"`
	Unet                 string     `json:"unet,omitempty"`
	Vae                  string     `json:"vae,omitempty"`
	VaeHash              string     `json:"vae_hash,omitempty"`
	Version              string     `json:"version,omitempty"`
	HiresResize          *Size      `json:"hires_resize,omitempty"`
}

// Float is a wrapper around float64 that supports unmarshaling from
// both numeric and string JSON values, including support for
// "inf", "Infinity", "-inf", "-Infinity" and "NaN" string values.
type Float float64

func (f *Float) UnmarshalJSON(v []byte) (err error) {
	var s string
	if err := json.Unmarshal(v, &s); err == nil {
		switch s {
		case "inf", "Infinity":
			*f = Float(math.Inf(1))
			return nil
		case "-inf", "-Infinity":
			*f = Float(math.Inf(-1))
			return nil
		case "NaN":
			*f = Float(math.NaN())
			return nil
		default:
			// float value as string
			var fv float64
			if fv, err = strconv.ParseFloat(s, 64); err != nil {
				return err
			}

			*f = Float(fv)

			return nil
		}
	}

	// just a regular float value
	var fv float64
	if err := json.Unmarshal(v, &fv); err != nil {
		return err
	}

	*f = Float(fv)

	return nil
}

func (f Float) MarshalJSON() ([]byte, error) {
	v := float64(f)
	switch {
	case math.IsNaN(v):
		return json.Marshal("NaN")
	case math.IsInf(v, 1):
		return json.Marshal("inf")
	case math.IsInf(v, -1):
		return json.Marshal("-inf")
	}

	return json.Marshal(v) // marshal result as standard float64
}

type FloatSlice []float32

func (s *FloatSlice) UnmarshalJSON(p []byte) error {
	if len(p) == 0 {
		return nil
	}

	if p[0] == '"' {
		var raw string
		if err := json.Unmarshal(p, &raw); err != nil {
			return err
		}

		parts := splitBracketList(raw)

		values := make(FloatSlice, 0, len(parts))
		for _, part := range parts {
			value, err := strconv.ParseFloat(part, 32)
			if err != nil {
				return err
			}

			values = append(values, float32(value))
		}

		*s = values

		return nil
	}

	var values []float32
	if err := json.Unmarshal(p, &values); err != nil {
		return err
	}

	*s = values

	return nil
}

type IntSlice []int

func (s *IntSlice) UnmarshalJSON(p []byte) error {
	if len(p) == 0 {
		return nil
	}

	if p[0] == '"' {
		var raw string
		if err := json.Unmarshal(p, &raw); err != nil {
			return err
		}

		parts := splitBracketList(raw)

		values := make(IntSlice, 0, len(parts))
		for _, part := range parts {
			value, err := strconv.Atoi(part)
			if err != nil {
				return err
			}

			values = append(values, value)
		}

		*s = values

		return nil
	}

	var values []int
	if err := json.Unmarshal(p, &values); err != nil {
		return err
	}

	*s = values

	return nil
}

type Loras []Lora

func (l *Loras) UnmarshalJSON(p []byte) (err error) {
	// Loras is a comma-separated list of Lora
	// Unmarshal to string first
	var tmp string
	if err := json.Unmarshal(p, &tmp); err != nil {
		return err
	}

	parts := strings.SplitSeq(tmp, ", ")
	for part := range parts {
		var lora Lora

		partString := fmt.Sprintf(`"%s"`, part)
		if err := json.Unmarshal([]byte(partString), &lora); err != nil {
			return err
		}

		*l = append(*l, lora)
	}

	return nil
}

type Lora struct {
	Name   string  `json:"name"`
	Weight float32 `json:"weight"`
}

func (s *Lora) UnmarshalJSON(p []byte) (err error) {
	// Lora is in the format "<lora:name:weight>"
	// Unmarshal to string first
	var tmp string
	if err := json.Unmarshal(p, &tmp); err != nil {
		return err
	}

	// Remove <lora: and >
	tmp = strings.TrimPrefix(tmp, "<lora:")
	tmp = strings.TrimSuffix(tmp, ">")

	parts := strings.SplitN(tmp, ":", 2)

	s.Name = parts[0]
	if len(parts) > 1 {
		weight, err := strconv.ParseFloat(parts[1], 32)
		if err != nil {
			return err
		}

		s.Weight = float32(weight)
	} else {
		s.Weight = 1.0
	}

	return nil
}

type Size struct {
	Width  int
	Height int
}

func (s *Size) UnmarshalJSON(p []byte) (err error) {
	// Size is in the format "512x512"
	// Unmarshal to string first
	var tmp string
	if err := json.Unmarshal(p, &tmp); err != nil {
		return err
	}

	size := strings.SplitN(tmp, "x", 2)

	if s.Width, err = strconv.Atoi(size[0]); err != nil {
		return err
	}

	if s.Height, err = strconv.Atoi(size[1]); err != nil {
		return err
	}

	return nil
}

func (s Size) MarshalJSON() ([]byte, error) {
	val := fmt.Sprintf("%dx%d", s.Width, s.Height)
	return json.Marshal(val)
}

func splitBracketList(raw string) []string {
	trimmed := strings.TrimSpace(raw)
	trimmed = strings.TrimPrefix(trimmed, "[")

	trimmed = strings.TrimSuffix(trimmed, "]")
	if trimmed == "" {
		return nil
	}

	parts := strings.Split(trimmed, ",")

	compact := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			compact = append(compact, part)
		}
	}

	return compact
}

func ParseParameters(in string) (meta Metadata, err error) {
	if json.Valid([]byte(in)) {
		return meta, fmt.Errorf("input is JSON, not plaintext")
	}

	// Parse a1111 parameters string; here be dragons
	kv := make(map[string]string)

	// This does not preserve newlines in the prompts!
	in2 := strings.ReplaceAll(in, "\n", ",")

	r := regexp.MustCompile("([^:,]+): ([^,]+)")

	matches := r.FindAllStringSubmatchIndex(in2, -1)

	var pmKey string
	var pmIdx int

	for i, match := range matches {
		// Match m[1]: the key
		m1 := in2[match[2]:match[3]]
		// Match m[2]: the value
		m2 := in2[match[4]:match[5]]

		// The first match is special: everything unmatched prior is the prompt
		if i == 0 && match[0] > 0 {
			prompt := in2[:match[0]-1]
			kv["prompt"] = strings.TrimSpace(prompt)
		}

		// If prev was negative prompt, the match was not sufficient,
		// fix it up
		if needsExtendedValue(pmKey) {
			kv[pmKey] = strings.TrimSpace(in2[pmIdx : match[0]-1])
		}

		// Normalize key: Lowercase, trim spaces, replace spaces with underscores
		// e.g. "Model hash" -> "model_hash"
		k := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(m1)), " ", "_")

		// Normalize value: Trim spaces
		v := strings.TrimSpace(m2)

		if val, ok := kv[k]; !ok {
			kv[k] = v
		} else {
			kv[k] = strings.Join([]string{val, v}, ", ")
		}

		// Remember previous key and m[2] index (value)
		pmKey = k
		pmIdx = match[4]
	}

	if needsExtendedValue(pmKey) {
		kv[pmKey] = strings.TrimSpace(in2[pmIdx:])
	}

	normalizeStableDiffusionCPPFields(kv)

	// LoRAs
	lr := regexp.MustCompile("<lora:[^>]+>")

	loraMatches := lr.FindAllString(kv["prompt"], -1)
	if len(loraMatches) > 0 {
		kv["loras"] = strings.Join(loraMatches, ", ")
	}

	kvByte, _ := json.Marshal(kv)
	err = json.Unmarshal(kvByte, &meta)

	return meta, err
}

func needsExtendedValue(key string) bool {
	switch key {
	case "negative_prompt", "custom_sigmas", "skip_layers":
		return true
	default:
		return false
	}
}

func normalizeStableDiffusionCPPFields(kv map[string]string) {
	value, ok := kv["hires_upscale"]
	if !ok {
		return
	}

	if _, err := strconv.ParseFloat(value, 32); err == nil {
		return
	}

	kv["hires_upscaler"] = value
	delete(kv, "hires_upscale")

	if scale, ok := kv["hires_scale"]; ok {
		kv["hires_upscale"] = scale
	}
}
