package fooocus

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Metadata in current format (Fooocus v2.2 and newer)
const metaV23Json = `{
	"adm_guidance": "(1.5, 0.8, 0.3)",
	"base_model": "juggernautXL_v8Rundiffusion",
	"base_model_hash": "aeb7e9e689",
	"clip_skip": 2,
	"full_negative_prompt": ["(worst quality, low quality, normal quality, lowres, low details, oversaturated, undersaturated, overexposed, underexposed, grayscale, bw, bad photo, bad photography, bad art:1.4), (watermark, signature, text font, username, error, logo, words, letters, digits, autograph, trademark, name:1.2), (blur, blurry, grainy), morbid, ugly, asymmetrical, mutated malformed, mutilated, poorly lit, bad shadow, draft, cropped, out of frame, cut off, censored, jpeg artifacts, out of focus, glitch, duplicate, (airbrushed, cartoon, anime, semi-realistic, cgi, render, blender, digital art, manga, amateur:1.3), (3D ,3D Game, 3D Game Scene, 3D Character:1.1), (bad hands, bad anatomy, bad body, bad face, bad teeth, bad arms, bad legs, deformities:1.3)", "anime, cartoon, graphic, (blur, blurry, bokeh), text, painting, crayon, graphite, abstract, glitch, deformed, mutated, ugly, disfigured"],
	"full_prompt": ["cinematic still A sunflower field . emotional, harmonious, vignette, 4k epic detailed, shot on kodak, 35mm photo, sharp focus, high budget, cinemascope, moody, epic, gorgeous, film grain, grainy", "A sunflower field, highly detailed, magic, peaceful, flowing, beautiful, atmosphere, radiant, magical, sharp focus, very coherent, intricate, elegant, epic, colorful, amazing composition, cinematic, artistic, fine detail, professional, clear, joyful, unique, expressive, cute, iconic, best, vivid, awesome, perfect, ambient background, pristine, creative"],
	"guidance_scale": 4,
	"lora_combined_1": "sd_xl_offset_example-lora_1.0 : 0.1",
	"loras": [["sd_xl_offset_example-lora_1.0", 0.1, "4852686128"]],
	"metadata_scheme": "fooocus",
	"negative_prompt": "",
	"performance": "Speed",
	"prompt": "A sunflower field",
	"prompt_expansion": "A sunflower field, highly detailed, magic, peaceful, flowing, beautiful, atmosphere, radiant, magical, sharp focus, very coherent, intricate, elegant, epic, colorful, amazing composition, cinematic, artistic, fine detail, professional, clear, joyful, unique, expressive, cute, iconic, best, vivid, awesome, perfect, ambient background, pristine, creative",
	"refiner_model": "None",
	"refiner_switch": 0.5,
	"resolution": "(512, 512)",
	"sampler": "dpmpp_2m_sde_gpu",
	"scheduler": "karras",
	"seed": "127589946317439009",
	"sharpness": 2,
	"steps": 30,
	"styles": "['Fooocus V2', 'Fooocus Enhance', 'Fooocus Sharp']",
	"vae": "Default (model)",
	"version": "Fooocus v2.5.5"
}`

var metaV23 = &MetadataV23{
	AdmGuidance:        AdmGuidanceOf(1.5, 0.8, 0.3),
	BaseModel:          "juggernautXL_v8Rundiffusion",
	BaseModelHash:      "aeb7e9e689",
	ClipSkip:           2,
	FullNegativePrompt: []string{"(worst quality, low quality, normal quality, lowres, low details, oversaturated, undersaturated, overexposed, underexposed, grayscale, bw, bad photo, bad photography, bad art:1.4), (watermark, signature, text font, username, error, logo, words, letters, digits, autograph, trademark, name:1.2), (blur, blurry, grainy), morbid, ugly, asymmetrical, mutated malformed, mutilated, poorly lit, bad shadow, draft, cropped, out of frame, cut off, censored, jpeg artifacts, out of focus, glitch, duplicate, (airbrushed, cartoon, anime, semi-realistic, cgi, render, blender, digital art, manga, amateur:1.3), (3D ,3D Game, 3D Game Scene, 3D Character:1.1), (bad hands, bad anatomy, bad body, bad face, bad teeth, bad arms, bad legs, deformities:1.3)", "anime, cartoon, graphic, (blur, blurry, bokeh), text, painting, crayon, graphite, abstract, glitch, deformed, mutated, ugly, disfigured"},
	FullPrompt:         []string{"cinematic still A sunflower field . emotional, harmonious, vignette, 4k epic detailed, shot on kodak, 35mm photo, sharp focus, high budget, cinemascope, moody, epic, gorgeous, film grain, grainy", "A sunflower field, highly detailed, magic, peaceful, flowing, beautiful, atmosphere, radiant, magical, sharp focus, very coherent, intricate, elegant, epic, colorful, amazing composition, cinematic, artistic, fine detail, professional, clear, joyful, unique, expressive, cute, iconic, best, vivid, awesome, perfect, ambient background, pristine, creative"},
	GuidanceScale:      4,
	Loras: []Lora{{
		Name:   "sd_xl_offset_example-lora_1.0",
		Weight: 0.1,
		Hash:   "4852686128",
	}},
	LoraCombined1:   &LoraCombined{Name: "sd_xl_offset_example-lora_1.0", Weight: 0.1},
	MetadataScheme:  "fooocus",
	NegativePrompt:  "",
	Performance:     "Speed",
	Prompt:          "A sunflower field",
	PromptExpansion: "A sunflower field, highly detailed, magic, peaceful, flowing, beautiful, atmosphere, radiant, magical, sharp focus, very coherent, intricate, elegant, epic, colorful, amazing composition, cinematic, artistic, fine detail, professional, clear, joyful, unique, expressive, cute, iconic, best, vivid, awesome, perfect, ambient background, pristine, creative",
	RefinerModel:    "None",
	RefinerSwitch:   0.5,
	Resolution:      ResolutionOf(512, 512),
	Sampler:         "dpmpp_2m_sde_gpu",
	Scheduler:       "karras",
	Seed:            "127589946317439009",
	Sharpness:       2,
	Steps:           30,
	Styles:          []string{"Fooocus V2", "Fooocus Enhance", "Fooocus Sharp"},
	Vae:             "Default (model)",
	Version:         "Fooocus v2.5.5",
}

const metaV23AltJson = `{
  "prompt": "A sunflower field",
  "negative_prompt": "",
  "prompt_expansion": "",
  "styles": "['Fooocus Enhance', 'Fooocus Sharp']",
  "performance": "Lightning",
  "steps": 4,
  "resolution": "(1152, 896)",
  "guidance_scale": 1.0,
  "sharpness": 0.0,
  "adm_guidance": "(1.0, 1.0, 0.0)",
  "base_model": "ponyDiffusionV6XL_v6TurboDPOMerge.safetensors",
  "refiner_model": "None",
  "refiner_switch": 1.0,
  "adaptive_cfg": 1.0,
  "clip_skip": 2,
  "sampler": "euler",
  "scheduler": "sgm_uniform",
  "vae": "Default (model)",
  "seed": "4781032431475889838",
  "lora_combined_1": "lora1.safetensors : 0.8",
  "lora_combined_2": "lora2.safetensors : 0.7",
  "lora_combined_3": "lora3.safetensors : 0.9",
  "lora_combined_4": "sdxl_lightning_4step_lora.safetensors : 1.0",
  "lora_combined_5": "sdxl_lightning_4step_lora.safetensors : 1.0",
  "metadata_scheme": false,
  "version": "Fooocus v2.5.0"
}`

var metaV23Alt = &MetadataV23{
	AdaptiveCfg:   1,
	AdmGuidance:   AdmGuidanceOf(1.0, 1.0, 0.0),
	BaseModel:     "ponyDiffusionV6XL_v6TurboDPOMerge.safetensors",
	BaseModelHash: "",
	ClipSkip:      2,
	GuidanceScale: 1,
	Loras: []Lora{{
		Name: "lora1.safetensors", Weight: 0.8,
	}, {
		Name: "lora2.safetensors", Weight: 0.7,
	}, {
		Name: "lora3.safetensors", Weight: 0.9,
	}, {
		Name: "sdxl_lightning_4step_lora.safetensors", Weight: 1.0,
	}, {
		Name: "sdxl_lightning_4step_lora.safetensors", Weight: 1.0,
	}},
	LoraCombined1:   &LoraCombined{Name: "lora1.safetensors", Weight: 0.8},
	LoraCombined2:   &LoraCombined{Name: "lora2.safetensors", Weight: 0.7},
	LoraCombined3:   &LoraCombined{Name: "lora3.safetensors", Weight: 0.9},
	LoraCombined4:   &LoraCombined{Name: "sdxl_lightning_4step_lora.safetensors", Weight: 1.0},
	LoraCombined5:   &LoraCombined{Name: "sdxl_lightning_4step_lora.safetensors", Weight: 1.0},
	MetadataScheme:  "fooocus",
	NegativePrompt:  "",
	Performance:     "Lightning",
	Prompt:          "A sunflower field",
	PromptExpansion: "",
	RefinerModel:    "None",
	RefinerSwitch:   1.0,
	Resolution:      ResolutionOf(1152, 896),
	Sampler:         "euler",
	Scheduler:       "sgm_uniform",
	Seed:            "4781032431475889838",
	Sharpness:       0,
	Steps:           4,
	Styles:          []string{"Fooocus Enhance", "Fooocus Sharp"},
	Vae:             "Default (model)",
	Version:         "Fooocus v2.5.0",
}

// Fooocus v2.2 private log metadata format
const metaV22Json = `{
	"prompt": "A sunflower field",
	"negative_prompt": "",
	"prompt_expansion": "A sunflower field, highly detailed, magic, peaceful, flowing, beautiful, atmosphere, radiant, magical, sharp focus, very coherent, intricate, elegant, epic, colorful, amazing composition, cinematic, artistic, fine detail, professional, clear, joyful, unique, expressive, cute, iconic, best, vivid, awesome, perfect, ambient background, pristine, creative",
	"styles": "['Fooocus V2', 'Fooocus Enhance', 'Fooocus Sharp', 'Fooocus Photograph', 'Fooocus Negative']",
	"performance": "Speed",
	"resolution": "(1152, 896)",
	"guidance_scale": 4,
	"sharpness": 2,
	"adm_guidance": "(1.5, 0.8, 0.3)",
	"base_model": "juggernautXL_v8Rundiffusion.safetensors",
	"refiner_model": "None",
	"refiner_switch": 0.5,
	"sampler": "dpmpp_2m_sde_gpu",
	"scheduler": "karras",
	"seed": 3380777627941884610,
	"lora_combined_1": "sd_xl_offset_example-lora_1.0.safetensors : 0.1",
	"lora_combined_2": "PetDinosaur-v2.safetensors : 0.8",
	"metadata_scheme": false,
	"version": "Fooocus v2.2.1"
}`

var metaV22 = &MetadataV22{
	MetadataV23: MetadataV23{
		AdmGuidance:     AdmGuidanceOf(1.5, 0.8, 0.3),
		BaseModel:       "juggernautXL_v8Rundiffusion.safetensors",
		BaseModelHash:   "",
		ClipSkip:        0,
		GuidanceScale:   4,
		LoraCombined1:   &LoraCombined{Name: "sd_xl_offset_example-lora_1.0.safetensors", Weight: 0.1},
		LoraCombined2:   &LoraCombined{Name: "PetDinosaur-v2.safetensors", Weight: 0.8},
		NegativePrompt:  "",
		Performance:     "Speed",
		Prompt:          "A sunflower field",
		PromptExpansion: "A sunflower field, highly detailed, magic, peaceful, flowing, beautiful, atmosphere, radiant, magical, sharp focus, very coherent, intricate, elegant, epic, colorful, amazing composition, cinematic, artistic, fine detail, professional, clear, joyful, unique, expressive, cute, iconic, best, vivid, awesome, perfect, ambient background, pristine, creative",
		RefinerModel:    "None",
		RefinerSwitch:   0.5,
		Resolution:      ResolutionOf(1152, 896),
		Sampler:         "dpmpp_2m_sde_gpu",
		Scheduler:       "karras",
		Sharpness:       2,
		Steps:           0,
		Styles:          []string{"Fooocus V2", "Fooocus Enhance", "Fooocus Sharp", "Fooocus Photograph", "Fooocus Negative"},
		Vae:             "",
		Version:         "Fooocus v2.2.1",
	},
	Seed:           3380777627941884610,
	MetadataScheme: false,
}

var metaV22Converted = &MetadataV23{
	AdmGuidance:   AdmGuidanceOf(1.5, 0.8, 0.3),
	BaseModel:     "juggernautXL_v8Rundiffusion.safetensors",
	BaseModelHash: "",
	ClipSkip:      0,
	GuidanceScale: 4,
	Loras: []Lora{{
		Name:   "sd_xl_offset_example-lora_1.0.safetensors",
		Weight: 0.1,
	}, {
		Name:   "PetDinosaur-v2.safetensors",
		Weight: 0.8,
	}},
	LoraCombined1:   &LoraCombined{Name: "sd_xl_offset_example-lora_1.0.safetensors", Weight: 0.1},
	LoraCombined2:   &LoraCombined{Name: "PetDinosaur-v2.safetensors", Weight: 0.8},
	MetadataScheme:  "fooocus",
	NegativePrompt:  "",
	Performance:     "Speed",
	Prompt:          "A sunflower field",
	PromptExpansion: "A sunflower field, highly detailed, magic, peaceful, flowing, beautiful, atmosphere, radiant, magical, sharp focus, very coherent, intricate, elegant, epic, colorful, amazing composition, cinematic, artistic, fine detail, professional, clear, joyful, unique, expressive, cute, iconic, best, vivid, awesome, perfect, ambient background, pristine, creative",
	RefinerModel:    "None",
	RefinerSwitch:   0.5,
	Resolution:      ResolutionOf(1152, 896),
	Sampler:         "dpmpp_2m_sde_gpu",
	Scheduler:       "karras",
	Seed:            "3380777627941884610",
	Sharpness:       2,
	Steps:           30,
	Styles:          []string{"Fooocus V2", "Fooocus Enhance", "Fooocus Sharp", "Fooocus Photograph", "Fooocus Negative"},
	Vae:             "",
	Version:         "Fooocus v2.2.1",
}

// Metadata in legacy format (Fooocus v2.1 and older)
const metaV21Json = `{
	"Prompt": "dinosaur leashed",
	"Negative Prompt": "",
	"Fooocus V2 Expansion": "dinosaur leashed, full color, cinematic, stunning, highly detailed, true colors, complex, elegant, symmetry, light, epic, great composition, creative, perfect, thought, best, real, novel, romantic, new, tender, cute, fancy, burning, nice, lovely, hopeful, pretty, artistic, surreal, inspiring, beautiful, dramatic, illuminated, amazing",
	"Styles": "['Fooocus V2', 'Dark Moody Atmosphere', 'Misc Horror']",
	"Performance": "Extreme Speed",
	"Resolution": "(1024, 1024)",
	"Sharpness": 0.0,
	"Guidance Scale": 1.0,
	"ADM Guidance": "(1.0, 1.0, 0.0)",
	"Base Model": "sd_xl_base_1.0_0.9vae.safetensors",
	"Refiner Model": "None",
	"Refiner Switch": 1.0,
	"Sampler": "lcm",
	"Scheduler": "lcm",
	"Seed": 6515525224486854496,
	"LoRA 1": "PetDinosaur-v2.safetensors : 1.65",
	"LoRA 6": "sdxl_lcm_lora.safetensors : 1.0",
	"Version": "v2.1.860"
}`

var metaV21 = &MetadataV21{
	AdmGuidance:        AdmGuidanceOf(1.0, 1.0, 0.0),
	BaseModel:          "sd_xl_base_1.0_0.9vae.safetensors",
	FooocusV2Expansion: "dinosaur leashed, full color, cinematic, stunning, highly detailed, true colors, complex, elegant, symmetry, light, epic, great composition, creative, perfect, thought, best, real, novel, romantic, new, tender, cute, fancy, burning, nice, lovely, hopeful, pretty, artistic, surreal, inspiring, beautiful, dramatic, illuminated, amazing",
	GuidanceScale:      1.0,
	Lora1:              &LoraCombined{Name: "PetDinosaur-v2.safetensors", Weight: 1.65},
	Lora6:              &LoraCombined{Name: "sdxl_lcm_lora.safetensors", Weight: 1.0},
	NegativePrompt:     "",
	Performance:        "Extreme Speed",
	Prompt:             "dinosaur leashed",
	RefinerModel:       "None",
	RefinerSwitch:      1.0,
	Resolution:         ResolutionOf(1024, 1024),
	Sampler:            "lcm",
	Scheduler:          "lcm",
	Seed:               6515525224486854496,
	Sharpness:          0.0,
	Styles:             []string{"Fooocus V2", "Dark Moody Atmosphere", "Misc Horror"},
	Version:            "v2.1.860",
}

var metaV21Converted = &MetadataV23{
	AdmGuidance:   AdmGuidanceOf(1.0, 1.0, 0.0),
	BaseModel:     "sd_xl_base_1.0_0.9vae.safetensors",
	GuidanceScale: 1.0,
	LoraCombined1: &LoraCombined{Name: "PetDinosaur-v2.safetensors", Weight: 1.65},
	Loras: []Lora{{
		Name:   "PetDinosaur-v2.safetensors",
		Weight: 1.65,
	}, {
		Name:   "sdxl_lcm_lora.safetensors",
		Weight: 1.0,
	}},
	MetadataScheme:  "fooocus",
	NegativePrompt:  "",
	Performance:     "Extreme Speed",
	Prompt:          "dinosaur leashed",
	PromptExpansion: "dinosaur leashed, full color, cinematic, stunning, highly detailed, true colors, complex, elegant, symmetry, light, epic, great composition, creative, perfect, thought, best, real, novel, romantic, new, tender, cute, fancy, burning, nice, lovely, hopeful, pretty, artistic, surreal, inspiring, beautiful, dramatic, illuminated, amazing",
	RefinerModel:    "None",
	RefinerSwitch:   1.0,
	Resolution:      ResolutionOf(1024, 1024),
	Sampler:         "lcm",
	Scheduler:       "lcm",
	Seed:            "6515525224486854496",
	Sharpness:       0.0,
	Steps:           8,
	Styles:          []string{"Fooocus V2", "Dark Moody Atmosphere", "Misc Horror"},
	Version:         "v2.1.860",
}

func TestDecodeMetadata_V21(t *testing.T) {
	var out *MetadataV21
	err := json.Unmarshal([]byte(metaV21Json), &out)
	require.NoError(t, err)
	assert.Equal(t, metaV21, out)
}

func TestEncodeMetadata_V21(t *testing.T) {
	encoded, err := json.Marshal(metaV21)
	require.NoError(t, err)

	t.Skip("Differing precision for floating point string representation")
	assert.JSONEq(t, metaV21Json, string(encoded))
}

func TestDecodeMetadata_V22(t *testing.T) {
	var decoded *MetadataV23
	err := json.Unmarshal([]byte(metaV22Json), &decoded)
	require.NoError(t, err)
	assert.Equal(t, metaV22Converted, decoded)
}

func TestEncodeMetadata_V22(t *testing.T) {
	encoded, err := json.Marshal(metaV22)
	require.NoError(t, err)

	t.Skip("v22 should omit some fields that are still included here")
	assert.JSONEq(t, metaV22Json, string(encoded))
}

func TestDecodeMetadata_V23(t *testing.T) {
	var decoded *MetadataV23
	err := json.Unmarshal([]byte(metaV23Json), &decoded)
	require.NoError(t, err)
	assert.Equal(t, metaV23, decoded)
}

func TestEncodeMetadata_V23(t *testing.T) {
	encoded, err := json.Marshal(metaV23)
	require.NoError(t, err)
	assert.JSONEq(t, metaV23Json, string(encoded))
}

func TestDecodeMetadata_V23_Alt(t *testing.T) {
	var decoded *MetadataV23
	err := json.Unmarshal([]byte(metaV23AltJson), &decoded)
	require.NoError(t, err)
	assert.Equal(t, metaV23Alt, decoded)
}

func TestEncodeMetadata_V23_Alt(t *testing.T) {
	encoded, err := json.Marshal(metaV23Alt)
	require.NoError(t, err)

	t.Skip("Only supports v23")
	assert.JSONEq(t, metaV23AltJson, string(encoded))
}

func TestConvertMetadata_V21(t *testing.T) {
	meta := ConvertV21ToV23(metaV21)
	assert.Equal(t, *metaV21Converted, meta)
}

func TestConvertMetadata_V22(t *testing.T) {
	meta := ConvertV22ToV23(metaV22)
	assert.Equal(t, *metaV22Converted, meta)
}

func TestDecodeMetadataAny_V21(t *testing.T) {
	var out *metadataAny
	err := json.Unmarshal([]byte(metaV21Json), &out)
	require.NoError(t, err)
	assert.Equal(t, metaV21, out.MetadataV21)
	assert.Equal(t, metaV21Converted, out.asMetadataV23())
}

func TestDecodeMetadataAny_V22(t *testing.T) {
	var out *metadataAny
	err := json.Unmarshal([]byte(metaV22Json), &out)
	require.NoError(t, err)
	assert.NotNil(t, out.MetadataV22)
	assert.Equal(t, metaV22Converted, out.asMetadataV23())
}

func TestDecodeMetadataAny_V23(t *testing.T) {
	var out *metadataAny
	err := json.Unmarshal([]byte(metaV23Json), &out)
	require.NoError(t, err)
	assert.NotNil(t, out.MetadataV23)
	assert.Equal(t, metaV23, out.MetadataV23)
	assert.Equal(t, metaV23, out.asMetadataV23())
}

func TestDecodeMetadataAny_V23_Alt(t *testing.T) {
	var out *metadataAny
	err := json.Unmarshal([]byte(metaV23AltJson), &out)
	require.NoError(t, err)
	assert.NotNil(t, out.MetadataV23)
	assert.Equal(t, metaV23Alt, out.MetadataV23)
	assert.Equal(t, metaV23Alt, out.asMetadataV23())
}

func TestEncodeMetadataAny_V21(t *testing.T) {
	t.Skip("Marshalling via metadataAny is not implemented")
	assert.Fail(t, "TODO")
}

func TestEncodeMetadataAny_V22(t *testing.T) {
	t.Skip("Marshalling via metadataAny is not implemented")
	assert.Fail(t, "TODO")
}

func TestEncodeMetadataAny_V23(t *testing.T) {
	meta := &metadataAny{
		MetadataV23: metaV23,
	}
	encoded, err := json.Marshal(meta)
	require.NoError(t, err)

	t.Skip("Marshalling via metadataAny is not implemented")
	assert.JSONEq(t, metaV23Json, string(encoded))
}
