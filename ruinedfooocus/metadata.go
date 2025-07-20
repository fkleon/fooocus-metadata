package ruinedfooocus

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const (
	Software = "RuinedFooocus"
)

type Metadata struct {
	BaseModel      string  `json:"base_model_name"`
	BaseModelHash  string  `json:"base_model_hash"`
	CfgScale       float32 `json:"cfg"`
	ClipSkip       uint8   `json:"clip_skip"`
	Denoise        any     `json:"denoise"`
	Height         uint16  `json:"height"`
	Loras          []Lora  `json:"loras"`
	NegativePrompt string  `json:"Negative"`
	Prompt         string  `json:"Prompt"`
	Sampler        string  `json:"sampler_name"`
	Scheduler      string  `json:"scheduler"`
	Seed           int     `json:"seed"`
	StartStep      uint8   `json:"start_step"`
	Steps          uint8   `json:"steps"`
	Version        string  `json:"software"`
	Width          uint16  `json:"width"`
}

// Encoded as nested list of format:
// list [string, string] (lora hash, lora details)
type Lora struct {
	Name   string
	Weight float32
	Hash   string
}

func (l *Lora) UnmarshalJSON(p []byte) error {
	var tmp []json.RawMessage
	if err := json.Unmarshal(p, &tmp); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[0], &l.Hash); err != nil {
		return err
	}

	// String of format "<weight> - <name>"
	var details string
	if err := json.Unmarshal(tmp[1], &details); err != nil {
		return err
	}

	loraCombined := strings.SplitN(details, " - ", 2)

	weight, err := strconv.ParseFloat(loraCombined[0], 32)
	if err != nil {
		return err
	}
	l.Weight = float32(weight)

	l.Name = loraCombined[1]

	return nil
}

func (l *Lora) MarshalJSON() ([]byte, error) {
	// Build details string
	details := fmt.Sprintf("%g - %v", l.Weight, l.Name)
	return json.Marshal([]interface{}{l.Hash, details})
}

func parseMetadata(parameters string) (meta Metadata, err error) {

	// Parse metadata
	err = json.Unmarshal([]byte(parameters), &meta)
	if err != nil {
		return meta, fmt.Errorf("%s: failed to read parameters: %w", Software, err)
	}

	return
}
