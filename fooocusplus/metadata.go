// Package fooocusplus implements reading and writing [FooocusPlus] metadata
// (image generation parameters).
//
// [FooocusPlus]: https://github.com/DavidDragonsage/FooocusPlus
package fooocusplus

import (
	"encoding/json"
	"fmt"

	"github.com/fkleon/fooocus-metadata/fooocus"
)

const (
	Software = "FooocusPlus"
)

type Metadata struct {
	AdmGuidance        *fooocus.AdmGuidance `json:"ADM Guidance"`
	BackendEngine      string               `json:"Backend Engine"`
	BaseModel          string               `json:"Base Model"`
	BaseModelHash      string               `json:"Base Model Hash"`
	ClipSkip           uint8                `json:"CLIP Skip"`
	FooocusV2Expansion string               `json:"Fooocus V2 Expansion"`
	FullNegativePrompt []string             `json:"Full Negative Prompt"`
	FullPrompt         []string             `json:"Full Prompt"`
	GuidanceScale      float32              `json:"Guidance Scale"`
	Loras              []fooocus.Lora       `json:"LoRAs"`
	MetadataScheme     string               `json:"Metadata Scheme"`
	NegativePrompt     string               `json:"Negative Prompt"`
	Performance        string               `json:"Performance"`
	Prompt             string               `json:"Prompt"`
	RefinerModel       string               `json:"Refiner Model,omitempty"`
	RefinerModelHash   string               `json:"Refiner Model Hash,omitempty"`  // TODO
	RefinerSwapMethod  string               `json:"Refiner Swap Method,omitempty"` // TODO
	RefinerSwitch      float32              `json:"Refiner Switch"`
	Resolution         *fooocus.Resolution  `json:"Resolution"`
	Sampler            string               `json:"Sampler"`
	Scheduler          string               `json:"Scheduler"`
	Seed               string               `json:"Seed"`
	Sharpness          float32              `json:"Sharpness"`
	Steps              uint8                `json:"Steps"`
	Styles             fooocus.Styles       `json:"Styles"`
	StylesDefinition   string               `json:"styles_definition"`
	User               string               `json:"User"`
	Vae                string               `json:"VAE"`
	Version            string               `json:"Version"`
}

type MetadataPrivateLog struct {
	AdmGuidance        *fooocus.AdmGuidance `json:"adm_guidance"`
	BackendEngine      string               `json:"backend_engine"`
	BaseModel          string               `json:"base_model"`
	BaseModelHash      string               `json:"Base Model Hash"` // TODO
	ClipSkip           uint8                `json:"clip_skip"`
	FooocusV2Expansion string               `json:"prompt_expansion"`
	FullNegativePrompt []string             `json:"Full Negative Prompt"` // TODO
	FullPrompt         []string             `json:"Full Prompt"`          // TODO
	GuidanceScale      float32              `json:"guidance_scale"`
	Loras              []fooocus.Lora       `json:"LoRAs"` // TODO
	MetadataScheme     string               `json:"metadata_scheme"`
	NegativePrompt     string               `json:"negative_prompt"`
	Performance        string               `json:"performance"`
	Prompt             string               `json:"prompt"`
	RefinerModel       string               `json:"refiner_model,omitempty"`
	RefinerModelHash   string               `json:"Refiner Model Hash,omitempty"`  // TODO
	RefinerSwapMethod  string               `json:"Refiner Swap Method,omitempty"` // TODO
	RefinerSwitch      float32              `json:"refiner_switch"`
	Resolution         *fooocus.Resolution  `json:"resolution"`
	Sampler            string               `json:"sampler"`
	Scheduler          string               `json:"scheduler"`
	Seed               string               `json:"seed"`
	Sharpness          float32              `json:"sharpness"`
	Steps              uint8                `json:"steps"`
	Styles             fooocus.Styles       `json:"styles"`
	StylesDefinition   string               `json:"styles_definition"` //TODO
	User               string               `json:"user,omitempty"`    // TODO
	Vae                string               `json:"vae"`
	Version            string               `json:"version"`
}

func (legacy *MetadataPrivateLog) toMetadata() (meta Metadata) {
	// Unsupported:
	// - BaseModelHash
	// - User
	// - FullNegativePrompt
	// - FullPrompt
	// - RefinerModelHash
	// - RefinerSwapMethod
	meta = Metadata{
		AdmGuidance:        legacy.AdmGuidance,
		BackendEngine:      legacy.BackendEngine,
		BaseModel:          legacy.BaseModel,
		ClipSkip:           legacy.ClipSkip,
		FooocusV2Expansion: legacy.FooocusV2Expansion,
		GuidanceScale:      legacy.GuidanceScale,
		//Loras:                {}, // TODO: post process default from LoraCombined fields
		MetadataScheme:    legacy.MetadataScheme,
		NegativePrompt:    legacy.NegativePrompt,
		Performance:       legacy.Performance,
		Prompt:            legacy.Prompt,
		RefinerModel:      legacy.RefinerModel,
		RefinerSwapMethod: legacy.RefinerSwapMethod,
		RefinerSwitch:     legacy.RefinerSwitch,
		Resolution:        legacy.Resolution,
		Sampler:           legacy.Sampler,
		Scheduler:         legacy.Scheduler,
		Seed:              legacy.Seed,
		Sharpness:         legacy.Sharpness,
		Steps:             legacy.Steps,
		Styles:            legacy.Styles,
		Vae:               legacy.Vae,
		Version:           legacy.Version,
	}
	return meta
}

func parseMetadata(parameters string) (meta Metadata, err error) {

	// Parse metadata
	err = json.Unmarshal([]byte(parameters), &meta)
	if err != nil {
		return meta, fmt.Errorf("%s: failed to read parameters: %w", Software, err)
	}

	return
}
