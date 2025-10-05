// Package stablediffusion implements reading [AUTOMATIC1111] style plaintext metadata
// (image generation parameters).
//
// [AUTOMATIC1111]: https://github.com/AUTOMATIC1111/stable-diffusion-webui
package stablediffusion

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Metadata struct {
	BatchSize            int     `json:"batch_size,string,omitempty"`
	BatchPos             int     `json:"batch_pos,string,omitempty"`
	CfgScale             float32 `json:"cfg_scale,string"`
	ClipSkip             int     `json:"clip_skip,string,omitempty"`
	DenoisingStrength    float32 `json:"denoising_strength,string,omitempty"`
	Eta                  float32 `json:"eta,string,omitempty"`
	HiresSteps           int     `json:"hires_steps,string,omitempty"`
	HiresUpscale         float32 `json:"hires_upscale,string,omitempty"`
	HiresUpscaler        string  `json:"hires_upscaler,omitempty"`
	Guidance             float32 `json:"guidance,string,omitempty"`
	ImageNoiseMultiplier float32 `json:"image_noise_multiplier,string,omitempty"`
	Loras                Loras   `json:"loras,omitempty"`
	Model                string  `json:"model,omitempty"`
	ModelHash            string  `json:"model_hash,omitempty"`
	NegativePrompt       string  `json:"negative_prompt,omitempty"`
	Prompt               string  `json:"prompt"`
	Rng                  string  `json:"rng,omitempty"`
	Sampler              string  `json:"sampler"` // The Sampler field contains both the sampler and scheduler names
	Seed                 int     `json:"seed,string"`
	Size                 *Size   `json:"size,omitempty"`
	Steps                int     `json:"steps,string"`
	VaeHash              string  `json:"vae_hash,omitempty"`
	Version              string  `json:"version,omitempty"`
}

type Loras []Lora

func (l *Loras) UnmarshalJSON(p []byte) (err error) {
	// Loras is a comma-separated list of Lora
	// Unmarshal to string first
	var tmp string
	if err := json.Unmarshal(p, &tmp); err != nil {
		return err
	}

	parts := strings.Split(tmp, ", ")
	for _, part := range parts {
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
		if weight, err := strconv.ParseFloat(parts[1], 32); err != nil {
			return err
		} else {
			s.Weight = float32(weight)
		}
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

	var pm_key string
	var pm_idx int

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
		if pm_key == "negative_prompt" {
			nprompt := in2[pm_idx : match[0]-1]
			kv["negative_prompt"] = strings.TrimSpace(nprompt)
		}

		// Normalize key: Lowercase, trim spaces, replace spaces with underscores
		// e.g. "Model hash" -> "model_hash"
		k := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(m1)), " ", "_")

		// Normalize value: Trim spaces
		v := strings.TrimSpace(m2)

		if _, ok := kv[k]; !ok {
			kv[k] = v
		}

		// Remember previous key and m[2] index (value)
		pm_key = k
		pm_idx = match[4]
	}

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
