package fooocus

import (
	"encoding/json"
	"testing"

	"github.com/cozy/goexif2/exif"
	"github.com/cozy/goexif2/tiff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Metadata in current format (Fooocus v2.2 and newer)
const metaJson = `{
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

var meta = &FooocusMeta{
	AdmGuidance:        [3]float32{1.5, 0.8, 0.3},
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
	Resolution:      [2]uint16{512, 512},
	Sampler:         "dpmpp_2m_sde_gpu",
	Scheduler:       "karras",
	Seed:            "127589946317439009",
	Sharpness:       2,
	Steps:           30,
	Styles:          []string{"Fooocus V2", "Fooocus Enhance", "Fooocus Sharp"},
	Vae:             "Default (model)",
	Version:         "Fooocus v2.5.5",
}

// Metadata in legacy format (Fooocus v2.1 and older)
const metaLegacyJson = `{
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

var metaLegacy = &FooocusMetaLegacy{
	AdmGuidance:        [3]float32{1.0, 1.0, 0.0},
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
	Resolution:         [2]uint16{1024, 1024},
	Sampler:            "lcm",
	Scheduler:          "lcm",
	Seed:               6515525224486854496,
	Sharpness:          0.0,
	Styles:             []string{"Fooocus V2", "Dark Moody Atmosphere", "Misc Horror"},
	Version:            "v2.1.860",
	// TODO: Steps derived from Performance mode
}

var metaLegacyConverted = &FooocusMeta{
	AdmGuidance:   [3]float32{1.0, 1.0, 0.0},
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
	Resolution:      [2]uint16{1024, 1024},
	Sampler:         "lcm",
	Scheduler:       "lcm",
	Seed:            "6515525224486854496",
	Sharpness:       0.0,
	Styles:          []string{"Fooocus V2", "Dark Moody Atmosphere", "Misc Horror"},
	Version:         "v2.1.860",
}

func TestDecodeLegacyMetadata(t *testing.T) {
	var out *FooocusMetaLegacy
	err := json.Unmarshal([]byte(metaLegacyJson), &out)
	require.NoError(t, err)
	assert.Equal(t, metaLegacy, out)
}

func TestEncodeLegacyMetadata(t *testing.T) {
	encoded, err := json.Marshal(metaLegacy)
	require.NoError(t, err)

	t.Skip("Differing precision for floating point string representation")
	assert.JSONEq(t, metaLegacyJson, string(encoded))
}

func TestConvertLegacyMetadata(t *testing.T) {
	meta := metaLegacy.toFooocusMeta()
	assert.Equal(t, metaLegacyConverted, meta)
}

func TestDecodeMetadata(t *testing.T) {
	var decoded *FooocusMeta
	err := json.Unmarshal([]byte(metaJson), &decoded)
	require.NoError(t, err)
	assert.Equal(t, meta, decoded)
}

func TestEncodeMetadata(t *testing.T) {
	encoded, err := json.Marshal(meta)
	require.NoError(t, err)
	assert.JSONEq(t, metaJson, string(encoded))
}

func TestExtractMetadataFromPNG(t *testing.T) {
	pngData := map[string]string{
		"fooocus_scheme": fooocus,
		"parameters":     metaJson,
	}
	fooocusData, err := ExtractMetadataFromPngData(pngData)
	require.NoError(t, err)
	assert.Equal(t, meta, fooocusData)
}

func TestExtractMetadataFromExif(t *testing.T) {
	t.Skip("Construct EXIF is a bit tricky")

	var exifData *exif.Exif = &exif.Exif{}
	var fieldMap = map[uint16]exif.FieldName{
		0x927c: exif.MakerNote,
		0x9286: exif.UserComment,
	}
	exifData.LoadTags(&tiff.Dir{
		Tags: []*tiff.Tag{
			{
				Id:  0x927c,
				Val: []byte(fooocus),
			},
			{
				Id:  0x9286,
				Val: []byte(metaJson),
			},
		},
	}, fieldMap, false)

	fooocusData, err := ExtractMetadataFromExifData(exifData)
	require.NoError(t, err)
	assert.Equal(t, meta, fooocusData)
}
