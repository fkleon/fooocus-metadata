// TODO update docs
// Package fooocus implements reading and writing [Fooocus] metadata (image generation parameters).
//
// To read embedded metadata, use [NewImageInfo]:
//
//	path := "testdata/sample.jpg"
//	image, err := NewImageInfo(path)
//	fmt.Println(image.FooocusMetadata.Version) // prints "Fooocus v2.5.5"
//
// To read from the private log file, use [ParsePrivateLog]:
//
//	path := "testdata/log.html"
//	images, err := ParsePrivateLog(privateLogFile)
//	meta := images["fooocus-meta.jpeg"]
//	fmt.Println(image.FooocusMetadata.Version) // prints "Fooocus v2.5.5"
//
// To write metadata into a PNG, use [EmbedMetadataAsPngText]:
//
//	meta := &Metadata{}
//	target, err := os.OpenFile("out.png", os.O_CREATE|os.O_WRONLY, 0644)
//	err = EmbedMetadataAsPngText(nil, target, meta)
//
// [Fooocus]: https://github.com/lllyasviel/Fooocus
package fooocus

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/bep/imagemeta"
)

const (
	Software = "Fooocus"
)

// Fooocus suports encoding metadata with one of two schemes:
//   - the native JSON scheme
//   - the AUTOMATIC1111 plaintext format for compatibility with Stable Diffusion web UI
//
//go:generate stringer -linecomment -type=MetadataScheme
type MetadataScheme uint8

const (
	Fooocus MetadataScheme = iota // fooocus
	A1111                         // a1111
)

type Performance uint8

const (
	// Default step configuration for performance presets
	quality      Performance = 60
	speed        Performance = 30
	extremeSpeed Performance = 8
	lightning    Performance = 4
	hyperSD      Performance = 4
)

// Fooocus metadata scheme (json).
//
// Implemented in Fooocus v2.2.0 and newer.
//
// Reference implementation:
//   - [feat: add metadata to images]
//   - [Serialisation]
//   - [Deserialisation]
//
// [feat: add metadata to images]: https://github.com/lllyasviel/Fooocus/pull/1940
// [Serialisation]: https://github.com/lllyasviel/Fooocus/blob/v2.5.5/modules/async_worker.py#L337
// [Deserialisation]: https://github.com/lllyasviel/Fooocus/blob/v2.5.5/modules/meta_parser.py#L22
type Metadata struct {
	AdaptiveCfg          float32       `json:"adaptive_cfg,omitempty"`
	AdmGuidance          *AdmGuidance  `json:"adm_guidance"`
	BaseModel            string        `json:"base_model"`
	BaseModelHash        string        `json:"base_model_hash"`
	ClipSkip             uint8         `json:"clip_skip"`
	CreatedBy            string        `json:"created_by,omitempty"`
	FreeU                *FreeU        `json:"freeu,omitempty"` // string: python tuple (b1: float, b2: float, s1: float, s2: float)
	FullNegativePrompt   []string      `json:"full_negative_prompt,omitempty"`
	FullPrompt           []string      `json:"full_prompt,omitempty"`
	GuidanceScale        float32       `json:"guidance_scale"`
	ImageNumber          uint          `json:"image_number,omitempty"`
	InpaintEngineVersion string        `json:"inpaint_engine_version,omitempty"`
	InpaintMode          string        `json:"inpaint_method,omitempty"`
	LoraCombined1        *LoraCombined `json:"lora_combined_1,omitempty"`
	LoraCombined2        *LoraCombined `json:"lora_combined_2,omitempty"`
	LoraCombined3        *LoraCombined `json:"lora_combined_3,omitempty"`
	LoraCombined4        *LoraCombined `json:"lora_combined_4,omitempty"`
	LoraCombined5        *LoraCombined `json:"lora_combined_5,omitempty"`
	Loras                []Lora        `json:"loras"`
	MetadataScheme       string        `json:"metadata_scheme"`
	NegativePrompt       string        `json:"negative_prompt"`
	Performance          string        `json:"performance"`
	Prompt               string        `json:"prompt"`
	PromptExpansion      string        `json:"prompt_expansion"`
	RefinerModel         string        `json:"refiner_model,omitempty"`
	RefinerModelHash     string        `json:"refiner_model_hash,omitempty"`
	RefinerSwapMethod    string        `json:"refiner_swap_method,omitempty"`
	RefinerSwitch        float32       `json:"refiner_switch"`
	Resolution           *Resolution   `json:"resolution"`
	Sampler              string        `json:"sampler"`
	Scheduler            string        `json:"scheduler"`
	Seed                 string        `json:"seed"`
	Sharpness            float32       `json:"sharpness"`
	Steps                uint8         `json:"steps"`
	Styles               Styles        `json:"styles"`
	Vae                  string        `json:"vae"`
	Version              string        `json:"version"`
}

func (meta *Metadata) fillLoras() {
	var loras []Lora = make([]Lora, 0, 5)

	var addLora = func(lora *LoraCombined) {
		if lora != nil {
			loras = append(loras, Lora{
				Name:   lora.Name,
				Weight: lora.Weight,
				Hash:   lora.Hash,
			})
		}
	}
	addLora(meta.LoraCombined1)
	addLora(meta.LoraCombined2)
	addLora(meta.LoraCombined3)
	addLora(meta.LoraCombined4)
	addLora(meta.LoraCombined5)

	meta.Loras = loras
}

// Fooocus legacy metadata schema (json).
//
// This format is only found in the private log HTML file generated by Fooocus
// v2.1.0 and older.
//
// Reference implementation:
//   - [Serialisation]
//   - [Deserialisation]
//
// [Serialisation]: https://github.com/lllyasviel/Fooocus/blob/v2.5.5/modules/async_worker.py#L337
// [Deserialisation]: https://github.com/lllyasviel/Fooocus/blob/v2.5.5/modules/meta_parser.py#L22
type MetadataLegacy struct {
	AdmGuidance          *AdmGuidance  `json:"ADM Guidance"` // string: python tuple (p: float, n: float, e: float)
	BaseModel            string        `json:"Base Model"`
	CFGMimicking         float32       `json:"CFG Mimicking from TSNR,omitempty"`
	ClipSkip             uint8         `json:"CLIP Skip,omitempty"`
	FooocusV2Expansion   string        `json:"Fooocus V2 Expansion"`
	FreeU                *FreeU        `json:"FreeU,omitempty"` // string: python tuple (b1: float, b2: float, s1: float, s2: float)
	GuidanceScale        float32       `json:"Guidance Scale"`
	ImageNumber          uint          `json:"Image Number,omitempty"`
	InpaintEngineVersion string        `json:"Inpaint Engine Version,omitempty"`
	InpaintMode          string        `json:"Inpaint Mode,omitempty"`
	Lora1                *LoraCombined `json:"LoRA 1,omitempty"`
	Lora2                *LoraCombined `json:"LoRA 2,omitempty"`
	Lora3                *LoraCombined `json:"LoRA 3,omitempty"`
	Lora4                *LoraCombined `json:"LoRA 4,omitempty"`
	Lora5                *LoraCombined `json:"LoRA 5,omitempty"`
	Lora6                *LoraCombined `json:"LoRA 6,omitempty"`
	NegativePrompt       string        `json:"Negative Prompt"`
	OverWriteSwitch      float32       `json:"Overwrite Switch,omitempty"`
	Performance          string        `json:"Performance"`
	Prompt               string        `json:"Prompt"`
	RefinerModel         string        `json:"Refiner Model"`
	RefinerSwapMethod    string        `json:"Refiner Swap Method,omitempty"`
	RefinerSwitch        float32       `json:"Refiner Switch"`
	Resolution           *Resolution   `json:"Resolution"` // string: python tuple (width: int, height: int)
	Sampler              string        `json:"Sampler"`
	Scheduler            string        `json:"Scheduler"`
	Seed                 int           `json:"Seed"`
	Sharpness            float32       `json:"Sharpness"`
	Steps                uint8         `json:"Steps,omitempty"` // TODO: post process default from Performance
	Styles               Styles        `json:"Styles"`
	Vae                  string        `json:"VAE,omitempty"`
	Version              string        `json:"Version"`
}

// Convert legacy metadata to current version
func (legacy *MetadataLegacy) toCurrent() (meta Metadata) {
	var loras []Lora = make([]Lora, 0, 6)

	var addLora = func(lora *LoraCombined) {
		if lora != nil {
			loras = append(loras, Lora{
				Name:   lora.Name,
				Weight: lora.Weight,
				Hash:   lora.Hash,
			})
		}
	}
	addLora(legacy.Lora1)
	addLora(legacy.Lora2)
	addLora(legacy.Lora3)
	addLora(legacy.Lora4)
	addLora(legacy.Lora5)
	addLora(legacy.Lora6)

	// Unsupported:
	// - BaseModelHash
	// - CreatedBy
	// - FullNegativePrompt
	// - FullPrompt
	// - RefinerModelHash
	meta = Metadata{
		AdaptiveCfg:          legacy.CFGMimicking,
		AdmGuidance:          legacy.AdmGuidance,
		BaseModel:            legacy.BaseModel,
		ClipSkip:             legacy.ClipSkip,
		FreeU:                legacy.FreeU,
		GuidanceScale:        legacy.GuidanceScale,
		ImageNumber:          legacy.ImageNumber,
		InpaintEngineVersion: legacy.InpaintEngineVersion,
		InpaintMode:          legacy.InpaintMode,
		LoraCombined1:        legacy.Lora1,
		LoraCombined2:        legacy.Lora2,
		LoraCombined3:        legacy.Lora3,
		LoraCombined4:        legacy.Lora4,
		LoraCombined5:        legacy.Lora5,
		Loras:                loras, // TODO: post process default from LoraCombined fields
		MetadataScheme:       Fooocus.String(),
		NegativePrompt:       legacy.NegativePrompt,
		Prompt:               legacy.Prompt,
		PromptExpansion:      legacy.FooocusV2Expansion,
		Performance:          legacy.Performance,
		RefinerModel:         legacy.RefinerModel,
		RefinerSwapMethod:    legacy.RefinerSwapMethod,
		RefinerSwitch:        legacy.RefinerSwitch,
		Resolution:           legacy.Resolution,
		Sampler:              legacy.Sampler,
		Scheduler:            legacy.Scheduler,
		Seed:                 strconv.Itoa(legacy.Seed),
		Sharpness:            legacy.Sharpness,
		Steps:                legacy.Steps,
		Styles:               legacy.Styles,
		Vae:                  legacy.Vae,
		Version:              legacy.Version,
	}
	return meta
}

// Tuple is an arbitrary length tuple backed by a slice.
// It is used to marshal/unmarshal Fooocus' string-encoded tuples.
type Tuple[T uint16 | float32] struct {
	data []T
}

func NewTuple[T uint16 | float32](data ...T) Tuple[T] {
	return Tuple[T]{
		data,
	}
}

func (r *Tuple[T]) UnmarshalJSON(p []byte) error {
	// Rewrite String-encoded Python tuple as JSON array:
	// "(1024, 1024)" -> [1024, 1024]
	pc := slices.Concat([]byte("["), p[2:len(p)-2], []byte("]"))

	return json.Unmarshal(pc, &r.data)
}

func (r *Tuple[T]) MarshalJSON() ([]byte, error) {
	var values []string = make([]string, len(r.data))

	for i, item := range r.data {
		switch it := any(item).(type) {
		case float32:
			values[i] = fmt.Sprintf("%.1f", it)
		case uint16:
			values[i] = fmt.Sprintf("%d", it)
		default:
			values[i] = fmt.Sprintf("%v", it)
		}
	}

	val := fmt.Sprintf("(%s)", strings.Join(values, ", "))
	return json.Marshal(val)
}

type Resolution struct {
	Tuple[uint16]
}

func ResolutionOf(height uint16, width uint16) *Resolution {
	return &Resolution{
		NewTuple(height, width),
	}
}

type FreeU struct {
	Tuple[float32]
}

func FreeUOf(b1 float32, b2 float32, s1 float32, s2 float32) *FreeU {
	return &FreeU{
		NewTuple(b1, b2, s1, s2),
	}
}

type AdmGuidance struct {
	Tuple[float32]
}

func AdmGuidanceOf(p float32, n float32, e float32) *AdmGuidance {
	return &AdmGuidance{
		NewTuple(p, n, e),
	}
}

// Styles are encoded within a string using single-quoted values, e.g.:
// "['Fooocus V2', 'Fooocus Enhance', 'Fooocus Sharp']"
type Styles []string

func (s *Styles) UnmarshalJSON(p []byte) error {
	var tmp string
	if err := json.Unmarshal(p, &tmp); err != nil {
		return err
	}

	var cleanStyles = strings.ReplaceAll(tmp, "'", "\"")
	var styles []string

	if err := json.Unmarshal([]byte(cleanStyles), &styles); err != nil {
		return err
	}
	*s = styles

	return nil
}

func (s *Styles) MarshalJSON() ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("[")
	for idx, style := range *s {
		sb.WriteString("'")
		sb.WriteString(style)
		sb.WriteString("'")
		if idx < len(*s)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")

	return json.Marshal(sb.String())
}

// Encoded as nested list of format:
// list [string, float32, string] (lora name, lora weight, lora hash)
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
	if err := json.Unmarshal(tmp[0], &l.Name); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[1], &l.Weight); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[2], &l.Hash); err != nil {
		return err
	}
	return nil
}

func (l *Lora) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{l.Name, l.Weight, l.Hash})
}

// String of format "<name> : <weight>"
type LoraCombined Lora

func (l *LoraCombined) UnmarshalJSON(p []byte) error {
	var tmp string
	if err := json.Unmarshal(p, &tmp); err != nil {
		return err
	}

	loraCombined := strings.SplitN(tmp, " : ", 2)
	l.Name = loraCombined[0]

	if len(loraCombined) > 1 {
		weight, err := strconv.ParseFloat(loraCombined[1], 32)
		if err != nil {
			return err
		}
		l.Weight = float32(weight)
	}
	return nil
}

func (l *LoraCombined) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%v : %g", l.Name, l.Weight))
}

func ExtractMetadataFromPngData(pngData map[string]string) (meta Metadata, err error) {

	if scheme, ok := pngData["fooocus_scheme"]; ok {
		parameters := pngData["parameters"]
		return parseMetadata(scheme, parameters)
	} else {
		return meta, fmt.Errorf("Fooocus: PNG: Metadata not found")
	}
}

func ExtractMetadataFromExifData(tags *imagemeta.Tags) (meta Metadata, err error) {

	var softwareVersion, scheme, parameters string

	exifData := tags.EXIF()

	if software, ok := exifData["Software"]; !ok {
		return meta, fmt.Errorf("Fooocus: EXIF: Software not found")
	} else {
		softwareVersion = software.Value.(string)
	}

	if !strings.HasPrefix(softwareVersion, "Fooocus ") {
		return meta, fmt.Errorf("Fooocus: EXIF: Unsupported software: %s", softwareVersion)
	}

	// imagemeta uses label "MakerNoteApple" for any "MakerNote" type.
	if makerNote, ok := exifData["MakerNoteApple"]; !ok {
		return meta, fmt.Errorf("Fooocus: EXIF: MakerNote not found")
	} else {
		scheme = makerNote.Value.(string)
	}

	if userComment, ok := exifData["UserComment"]; !ok {
		return meta, fmt.Errorf("Fooocus: EXIF: UserComment not found")
	} else {
		parameters = userComment.Value.(string)
	}

	return parseMetadata(scheme, parameters)
}

func parseMetadata(scheme string, parameters string) (meta Metadata, err error) {

	// Scheme is one of 'fooocus' or 'a1111'
	if scheme != Fooocus.String() {
		return meta, fmt.Errorf("Fooocus: unsupported metadata scheme: %s", scheme)
	}

	// Parse metadata
	err = json.Unmarshal([]byte(parameters), &meta)
	if err != nil {
		return meta, fmt.Errorf("Fooocus: failed to read parameters: %w", err)
	}

	return
}
