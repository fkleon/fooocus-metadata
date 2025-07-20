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
	"log/slog"
	"slices"
	"strconv"
	"strings"
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

//go:generate stringer -type=MetadataVersion
type MetadataVersion uint8

const (
	v21     MetadataVersion = iota + 21 // v2.1 ("legacy")
	v22                                 // v2.2
	v23                                 // v2.3+ ("current")
	unknown = 0
)

var (
	perf = map[string]uint8{
		"Quality":       60,
		"Speed":         30,
		"Extreme Speed": 8,
		"Lightning":     4,
		"Hyper-SD":      4,
	}
)

type Version struct {
	Version        string          `json:"version"`
	MetadataScheme json.RawMessage `json:"metadata_scheme"`
}

func (v *Version) MetadataVersion() MetadataVersion {
	if strings.HasPrefix(v.Version, "v2.1") {
		return v21
	} else if strings.HasPrefix(v.Version, "Fooocus v2.2") {
		return v22
	} else if strings.HasPrefix(v.Version, "Fooocus v2.3") ||
		strings.HasPrefix(v.Version, "Fooocus v2.4") ||
		strings.HasPrefix(v.Version, "Fooocus v2.5") {
		return v23
	} else {
		return unknown
	}
}

type metadataAny struct {
	Version
	*MetadataV21 // Fooocus v2.1 metadata structure ("legacy")
	*MetadataV22 // Fooocus v2.2 metadata structure
	*MetadataV23 // Fooocus v2.3+ metadata structure ("current")
}

func (m *metadataAny) asMetadataV23() *MetadataV23 {
	if m.MetadataV23 != nil {
		return m.MetadataV23
	} else if m.MetadataV22 != nil {
		current := ConvertV22ToV23(m.MetadataV22)
		return &current
	} else if m.MetadataV21 != nil {
		current := ConvertV21ToV23(m.MetadataV21)
		return &current
	}
	return nil
}

func (m *metadataAny) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &m.Version); err != nil {
		return err
	}

	switch v := m.MetadataVersion(); v {
	case unknown:
		slog.Warn("Unknown Fooocus metadata version", "version", v)
		return fmt.Errorf("Unknown Foooocus metadata version: %s", v)
	case v21:
		m.MetadataV21 = &MetadataV21{}
		return json.Unmarshal(data, m.MetadataV21)
	case v22:
		m.MetadataV22 = &MetadataV22{}
		return json.Unmarshal(data, m.MetadataV22)
	case v23:
		fallthrough
	default:
		m.MetadataV23 = &MetadataV23{}
		return json.Unmarshal(data, m.MetadataV23)
	}
}

func (m *metadataAny) MarshallJSON() ([]byte, error) {
	switch {
	case m.MetadataV21 != nil:
		return json.Marshal(m.MetadataV21)
	case m.MetadataV22 != nil:
		return json.Marshal(m.MetadataV22)
	case m.MetadataV23 != nil:
		return json.Marshal(m.MetadataV23)
	default:
		return json.Marshal(nil)
	}
}

type Metadata = MetadataV23

// Fooocus v2.3 metadata scheme (json).
//
// Used by Fooocus v2.2 and newer for embedded metadata.
// Used by Fooocus v2.3 and newer for private log metadata.
//
// Reference implementation:
//   - [feat: add metadata to images]
//   - [Serialisation]
//   - [Deserialisation]
//
// [feat: add metadata to images]: https://github.com/lllyasviel/Fooocus/pull/1940
// [Serialisation]: https://github.com/lllyasviel/Fooocus/blob/v2.5.5/modules/async_worker.py#L337
// [Deserialisation]: https://github.com/lllyasviel/Fooocus/blob/v2.5.5/modules/meta_parser.py#L22
type MetadataV23 struct {
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

// Fooocus v2.2 metadata scheme (json).
//
// This format is found in the private log HTML file generated by Fooocus
// v2.2.x.
type MetadataV22 struct {
	MetadataV23
	Seed           int  `json:"seed"`
	MetadataScheme bool `json:"metadata_scheme"`
}

func (m *MetadataV23) UnmarshalJSON(data []byte) error {
	// A variant of the v23 metadata uses different datatypes
	// for seed and metadata scheme.

	// Temporary type without UnmarshalJSON to avoid infinite
	// recursion.
	type metadata MetadataV23

	var dest = struct {
		*metadata
		// Exploits the fact that JSON unmarshalls into the
		// field with the shallowest depth.
		Seed           json.Number     `json:"seed"`
		MetadataScheme json.RawMessage `json:"metadata_scheme"`
	}{
		metadata: (*metadata)(m),
	}

	if err := json.Unmarshal(data, &dest); err != nil {
		return err
	}

	// Convert and populate values on the v23 struct.
	m.Seed = dest.Seed.String()

	// MetadataScheme either contains the name of the schema
	// or a boolean indicating whether metadata was embedded.
	if err := json.Unmarshal(dest.MetadataScheme, &m.MetadataScheme); err != nil {
		// Default 'fooocus' schema if not embedded.
		m.MetadataScheme = Fooocus.String()
	}

	m.fillLoras()
	m.fillSteps()
	return nil
}

func (meta *MetadataV23) fillLoras() {
	// If Loras are already set, do not overwrite them
	if meta.Loras != nil {
		return
	}

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

func (meta *Metadata) fillSteps() {
	// Set default steps based on performance preset
	if meta.Steps == 0 {
		if steps, ok := perf[meta.Performance]; ok {
			meta.Steps = steps
		}
	}
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
type MetadataV21 struct {
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

func ConvertV21ToV23(v21 *MetadataV21) (v23 Metadata) {
	legacy := v21
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
	v23 = Metadata{
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
		Loras:                loras,
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

	// Populate missing steps from performance preset
	v23.fillSteps()

	return v23
}

func ConvertV22ToV23(v22 *MetadataV22) (v23 MetadataV23) {
	v23 = v22.MetadataV23
	if v23.Seed == "" {
		v23.Seed = strconv.Itoa(v22.Seed)
	}
	v23.MetadataScheme = Fooocus.String()

	// Populate missing steps and LoRAs
	v23.fillSteps()
	v23.fillLoras()

	return v23
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

func parseMetadata(scheme string, parameters string) (meta Metadata, err error) {

	// Scheme is one of 'fooocus' or 'a1111'
	if scheme != Fooocus.String() {
		return meta, fmt.Errorf("%s: unsupported metadata scheme: %s", Software, scheme)
	}

	// Parse metadata
	err = json.Unmarshal([]byte(parameters), &meta)
	if err != nil {
		return meta, fmt.Errorf("%s: failed to read parameters: %w", Software, err)
	}

	return
}
