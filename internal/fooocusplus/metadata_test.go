package fooocusplus

import (
	"encoding/json"
	"testing"

	"github.com/bep/imagemeta"
	"github.com/fkleon/fooocus-metadata/internal/fooocus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Metadata as JSON
const metaJson = `{
  "ADM Guidance": "(1.5, 0.8, 0.3)",
  "Backend Engine": "SDXL-Fooocus",
  "Base Model": "elsewhereXL_v10",
  "Base Model Hash": "79fd29ab43",
  "CLIP Skip": 2,
  "Fooocus V2 Expansion": "A sunflower field, beautiful dynamic dramatic bright shining atmosphere, intricate, cinematic, extremely aesthetic, highly detailed, transparent, sharp focus, professional, winning, best, creative, cute, attractive, innocent, enhanced, colorful, color, inspired, pretty, depicted, illuminated, epic, amazing, artistic, pure, rational, elite, futuristic, inspirational, coherent, fantastic",
  "Full Negative Prompt": [
    "(worst quality, low quality, normal quality, lowres, low details, oversaturated, undersaturated, overexposed, underexposed, grayscale, bw, bad photo, bad photography, bad art:1.4), (watermark, signature, text font, username, error, logo, words, letters, digits, autograph, trademark, name:1.2), (blur, blurry, grainy), morbid, ugly, asymmetrical, mutated malformed, mutilated, poorly lit, bad shadow, draft, cropped, out of frame, cut off, censored, jpeg artifacts, out of focus, glitch, duplicate, (airbrushed, cartoon, anime, semi-realistic, cgi, render, blender, digital art, manga, amateur:1.3), (3D ,3D Game, 3D Game Scene, 3D Character:1.1), (bad hands, bad anatomy, bad body, bad face, bad teeth, bad arms, bad legs, deformities:1.3)"
  ],
  "Full Prompt": [
    "A sunflower field",
    "A sunflower field, beautiful dynamic dramatic bright shining atmosphere, intricate, cinematic, extremely aesthetic, highly detailed, transparent, sharp focus, professional, winning, best, creative, cute, attractive, innocent, enhanced, colorful, color, inspired, pretty, depicted, illuminated, epic, amazing, artistic, pure, rational, elite, futuristic, inspirational, coherent, fantastic"
  ],
  "Guidance Scale": 4.5,
  "LoRAs": [],
  "Metadata Scheme": "Fooocus",
  "Negative Prompt": "",
  "Performance": "Speed",
  "Prompt": "A sunflower field",
  "Refiner Model": "None",
  "Refiner Switch": 0.6,
  "Resolution": "(1024, 1024)",
  "Sampler": "dpmpp_2m_sde_gpu",
  "Scheduler": "karras",
  "Seed": "5256010854089202552",
  "Sharpness": 6,
  "Steps": 30,
  "Styles": "['Fooocus V2', 'Fooocus Enhance']",
  "User": "FooocusPlus",
  "VAE": "Default (model)",
  "Version": "FooocusPlus 1.0.0",
  "styles_definition": ""
}`

var meta = &Metadata{
	AdmGuidance:        fooocus.AdmGuidanceOf(1.5, 0.8, 0.3),
	BackendEngine:      "SDXL-Fooocus",
	BaseModel:          "elsewhereXL_v10",
	BaseModelHash:      "79fd29ab43",
	ClipSkip:           2,
	FooocusV2Expansion: "A sunflower field, beautiful dynamic dramatic bright shining atmosphere, intricate, cinematic, extremely aesthetic, highly detailed, transparent, sharp focus, professional, winning, best, creative, cute, attractive, innocent, enhanced, colorful, color, inspired, pretty, depicted, illuminated, epic, amazing, artistic, pure, rational, elite, futuristic, inspirational, coherent, fantastic",
	FullNegativePrompt: []string{
		"(worst quality, low quality, normal quality, lowres, low details, oversaturated, undersaturated, overexposed, underexposed, grayscale, bw, bad photo, bad photography, bad art:1.4), (watermark, signature, text font, username, error, logo, words, letters, digits, autograph, trademark, name:1.2), (blur, blurry, grainy), morbid, ugly, asymmetrical, mutated malformed, mutilated, poorly lit, bad shadow, draft, cropped, out of frame, cut off, censored, jpeg artifacts, out of focus, glitch, duplicate, (airbrushed, cartoon, anime, semi-realistic, cgi, render, blender, digital art, manga, amateur:1.3), (3D ,3D Game, 3D Game Scene, 3D Character:1.1), (bad hands, bad anatomy, bad body, bad face, bad teeth, bad arms, bad legs, deformities:1.3)",
	},
	FullPrompt: []string{
		"A sunflower field",
		"A sunflower field, beautiful dynamic dramatic bright shining atmosphere, intricate, cinematic, extremely aesthetic, highly detailed, transparent, sharp focus, professional, winning, best, creative, cute, attractive, innocent, enhanced, colorful, color, inspired, pretty, depicted, illuminated, epic, amazing, artistic, pure, rational, elite, futuristic, inspirational, coherent, fantastic",
	},
	GuidanceScale:  4.5,
	Loras:          []fooocus.Lora{},
	MetadataScheme: "Fooocus",
	NegativePrompt: "",
	Performance:    "Speed",
	Prompt:         "A sunflower field",
	RefinerModel:   "None",
	RefinerSwitch:  0.6,
	Resolution:     fooocus.ResolutionOf(1024, 1024),
	Sampler:        "dpmpp_2m_sde_gpu",
	Scheduler:      "karras",
	Seed:           "5256010854089202552",
	Sharpness:      6,
	Steps:          30,
	Styles: fooocus.Styles{
		"Fooocus V2",
		"Fooocus Enhance",
	},
	User:             "FooocusPlus",
	Vae:              "Default (model)",
	Version:          "FooocusPlus 1.0.0",
	StylesDefinition: "",
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

func TestExtractMetadataFromPNG(t *testing.T) {
	pngData := map[string]string{
		"Comment": metaJson,
	}
	fooocusData, err := ExtractMetadataFromPngData(pngData)
	require.NoError(t, err)
	assert.Equal(t, *meta, fooocusData)
}

func TestExtractMetadataFromExif(t *testing.T) {

	var exifData imagemeta.Tags

	exifData.Add(imagemeta.TagInfo{
		Source: imagemeta.EXIF,
		Tag:    "Software",
		Value:  "FooocusPlus 1.0.0",
	})
	exifData.Add(imagemeta.TagInfo{
		Source: imagemeta.EXIF,
		Tag:    "UserComment",
		Value:  metaJson,
	})

	fooocusData, err := ExtractMetadataFromExifData(&exifData)
	require.NoError(t, err)
	assert.Equal(t, *meta, fooocusData)
}
