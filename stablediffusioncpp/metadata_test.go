package stablediffusioncpp

import (
	"testing"

	sd "github.com/fkleon/fooocus-metadata/stablediffusion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseParameters_FromMetadataField(t *testing.T) {
	in := "ignored prompt\nNegative prompt: ignored negative\nSteps: 1, CFG scale: 99, Seed: 1, Size: 64x64, Model: ignored.safetensors, RNG: cpu, Sampler: euler, Version: stable-diffusion.cpp, SDCPP: {\"schema\":\"sdcpp.image.params/v1\",\"generator\":{\"name\":\"stable-diffusion.cpp\",\"version\":\"1.0.0\",\"commit\":\"abc123\"},\"seed\":1375127038,\"width\":512,\"height\":512,\"prompt\":{\"positive\":\"A sunflower field\",\"negative\":\"Blue sky\"},\"sampling\":{\"steps\":30,\"eta\":0,\"method\":\"dpm++2mv2\",\"scheduler\":\"karras\",\"guidance\":{\"txt_cfg\":5,\"img_cfg\":5,\"distilled_guidance\":3.5,\"slg\":{\"scale\":0,\"layers\":[],\"start\":0.01,\"end\":0.2}}},\"models\":{\"model\":\"v1-5-pruned-emaonly.safetensors\",\"clip_l\":\"clip_l.safetensors\",\"t5xxl\":\"t5-v1_1-xxl.gguf\",\"vae\":\"vae.safetensors\"},\"clip_skip\":2,\"strength\":0.4,\"rng\":\"cuda\"}"

	meta, err := ParseParameters(in)
	require.NoError(t, err)
	assert.Equal(t, sd.Metadata{
		CfgScale:          5,
		ClipSkip:          2,
		DenoisingStrength: 0.4,
		Eta:               0,
		Guidance:          3.5,
		Model:             "v1-5-pruned-emaonly.safetensors",
		NegativePrompt:    "Blue sky",
		Prompt:            "A sunflower field",
		Rng:               "cuda",
		Sampler:           "dpm++2mv2 karras",
		Seed:              1375127038,
		Size:              &sd.Size{Width: 512, Height: 512},
		Steps:             30,
		TextEncoder:       "clip_l.safetensors, t5-v1_1-xxl.gguf",
		Vae:               "vae.safetensors",
		Version:           "stable-diffusion.cpp (1.0.0)",
	}, meta)
}

func TestParseParameters_RawJSON(t *testing.T) {
	in := `{"schema":"sdcpp.image.params/v1","generator":{"name":"stable-diffusion.cpp","version":"1.18.0","commit":"deadbeef"},"seed":2056899594,"width":512,"height":512,"prompt":{"positive":"Dog with short hair.","negative":""},"sampling":{"steps":8,"eta":null,"method":"euler_a","scheduler":"smoothstep","guidance":{"txt_cfg":1,"img_cfg":1,"distilled_guidance":3.5,"slg":{"scale":0,"layers":[],"start":0.01,"end":0.2}}},"models":{"llm":"Qwen3-4B-Instruct-2507-Q3_K_S.gguf","diffusion_model":"z_image_turbo-Q6_K.gguf","vae":"ae-f16.gguf"},"rng":"cuda"}`

	meta, err := ParseParameters(in)
	require.NoError(t, err)
	assert.Equal(t, sd.Metadata{
		CfgScale:    1,
		Eta:         0,
		Guidance:    3.5,
		Model:       "",
		Prompt:      "Dog with short hair.",
		Rng:         "cuda",
		Sampler:     "euler_a smoothstep",
		Seed:        2056899594,
		Size:        &sd.Size{Width: 512, Height: 512},
		Steps:       8,
		TextEncoder: "Qwen3-4B-Instruct-2507-Q3_K_S.gguf",
		Unet:        "z_image_turbo-Q6_K.gguf",
		Vae:         "ae-f16.gguf",
		Version:     "stable-diffusion.cpp (1.18.0)",
	}, meta)
}

func TestParseParameters_HiresAndLoras(t *testing.T) {
	in := `Steps: 20, SDCPP: {"schema":"sdcpp.image.params/v1","generator":{"name":"stable-diffusion.cpp","version":"1.18.0","commit":"deadbeef"},"seed":42,"width":1024,"height":1024,"prompt":{"positive":"score_9, sunflower","negative":"corn"},"sampling":{"steps":20,"eta":0,"method":"euler_a","scheduler":"karras","guidance":{"txt_cfg":4,"img_cfg":4,"distilled_guidance":3.5,"slg":{"scale":0,"layers":[],"start":0.01,"end":0.2}}},"models":{"model":"sdxl.safetensors","clip_l":"clip_l.safetensors","clip_g":"clip_g.safetensors","diffusion_model":"sdxl-unet.gguf","vae":"sdxl-vae-fp16-fix.safetensors"},"clip_skip":2,"strength":0.2,"rng":"cuda","loras":[{"name":"Agriculture_V1.safetensors","multiplier":1,"is_high_noise":false},{"name":"size_slider_v1.safetensors","multiplier":1.7,"is_high_noise":true}],"hires":{"enabled":true,"upscaler":"4x-UltraSharp","model":"","scale":2,"target_width":0,"target_height":0,"steps":30,"denoising_strength":0.45,"upscale_tile_size":128}}`

	meta, err := ParseParameters(in)
	require.NoError(t, err)
	assert.Equal(t, sd.Metadata{
		CfgScale:          4,
		ClipSkip:          2,
		DenoisingStrength: 0.45,
		Eta:               0,
		Guidance:          3.5,
		HiresSteps:        30,
		HiresUpscale:      2,
		HiresUpscaler:     "4x-UltraSharp",
		Loras:             sd.Loras{{Name: "Agriculture_V1.safetensors", Weight: 1}, {Name: "size_slider_v1.safetensors", Weight: 1.7}},
		Model:             "sdxl.safetensors",
		NegativePrompt:    "corn",
		Prompt:            "score_9, sunflower",
		Rng:               "cuda",
		Sampler:           "euler_a karras",
		Seed:              42,
		Size:              &sd.Size{Width: 1024, Height: 1024},
		Steps:             20,
		TextEncoder:       "clip_l.safetensors, clip_g.safetensors",
		Unet:              "sdxl-unet.gguf",
		Vae:               "sdxl-vae-fp16-fix.safetensors",
		Version:           "stable-diffusion.cpp (1.18.0)",
	}, meta)
}

func TestParseParameters_Error(t *testing.T) {
	t.Run("missing field", func(t *testing.T) {
		_, err := ParseParameters("Steps: 20, CFG scale: 7")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "SDCPP field not found")
	})

	t.Run("invalid json", func(t *testing.T) {
		_, err := ParseParameters("Steps: 20, SDCPP: {not-json}")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid SDCPP JSON")
	})

	t.Run("unsupported schema", func(t *testing.T) {
		_, err := ParseParameters(`{"schema":"other.schema/v1"}`)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported SDCPP schema")
	})
}
