package stablediffusion

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	in  string
	out Metadata
}

func TestDecodeMetadata_SD_WebUI(t *testing.T) {
	tc := []testCase{
		{
			// wiki hero image
			in: "Astronaut in a jungle, cold color palette, muted colors, detailed, 8k\nSteps: 50, Sampler: DPM++ 2M Karras, CFG scale: 5, Seed: 42, Size: 1024x1024, Model hash: 1f69731261, Model: sd_xl_base_0.9, Clip skip: 2, RNG: CPU, Version: v1.4.1-166-g21aec6f5",
			out: Metadata{
				CfgScale:  5,
				ClipSkip:  2,
				Model:     "sd_xl_base_0.9",
				ModelHash: "1f69731261",
				Prompt:    "Astronaut in a jungle, cold color palette, muted colors, detailed, 8k",
				Rng:       "CPU",
				Sampler:   "DPM++ 2M Karras",
				Seed:      42,
				Size:      &Size{Width: 1024, Height: 1024},
				Steps:     50,
				Version:   "v1.4.1-166-g21aec6f5",
			},
		},
		{
			// outpainting-2.png
			in: "clouds, sky, nebula, 8k clean, cinematic lighting, highly detailed, digital painting, clean 8k, art by Roy Liechtestein.\nSteps: 130, Sampler: Euler a, CFG scale: 15, Seed: 4051576822, Denoising Strength: 1",
			out: Metadata{
				Steps:             130,
				Sampler:           "Euler a",
				CfgScale:          15,
				Seed:              4051576822,
				DenoisingStrength: 1,
				Prompt:            "clouds, sky, nebula, 8k clean, cinematic lighting, highly detailed, digital painting, clean 8k, art by Roy Liechtestein.",
			},
		},
		{
			// inpainting-81-euler-a.png
			in: "8K clean, underwater distortion, ripples, sea, corals,, underwater, high quality, award winning, photo\nSteps: 81, Sampler: Euler a, CFG scale: 15, Seed: 952858003, Denoising Strength: 1",
			out: Metadata{
				CfgScale:          15,
				DenoisingStrength: 1,
				Prompt:            "8K clean, underwater distortion, ripples, sea, corals,, underwater, high quality, award winning, photo",
				Sampler:           "Euler a",
				Seed:              952858003,
				Steps:             81,
			},
		},
		{
			// inpaint-mask2.png
			in: "sci-fi treehouse, cyberpunk, neon lights, intricate, cinematic lighting, highly detailed, digital painting, artstation, concept art, smooth, sharp focus, illustration, art by Gareth Pugh, Alex Timmermans, Abraham Mintchine, Alson S. Clark\nSteps: 78, Sampler: Euler a, CFG scale: 12, Seed: 833664775, Denoising strength: 1, Denoising strength change factor: 1",
			out: Metadata{
				CfgScale:          12,
				DenoisingStrength: 1,
				//DenoisingStrengthChangeFactor: 1,
				Prompt:  "sci-fi treehouse, cyberpunk, neon lights, intricate, cinematic lighting, highly detailed, digital painting, artstation, concept art, smooth, sharp focus, illustration, art by Gareth Pugh, Alex Timmermans, Abraham Mintchine, Alson S. Clark",
				Sampler: "Euler a",
				Seed:    833664775,
				Steps:   78,
			},
		},
		{
			// hi-res fix
			in: "cornfield, modern style, detailed face, beautiful face, by greg rutkowski and alphonse mucha, d & d character, in front of an urban background, digital painting, concept art, smooth, sharp focus illustration, artstation hq\nNegative prompt: ((((mutated hands and fingers))))\nSteps: 20, Sampler: Euler a, CFG scale: 12, Seed: 950170121, Size: 960x960, Model hash: 6ecd8e48, Batch size: 2, Batch pos: 0, Denoising strength: 0.7",
			out: Metadata{
				BatchPos:          0,
				BatchSize:         2,
				CfgScale:          12,
				DenoisingStrength: 0.7,
				ModelHash:         "6ecd8e48",
				NegativePrompt:    "((((mutated hands and fingers))))",
				Prompt:            "cornfield, modern style, detailed face, beautiful face, by greg rutkowski and alphonse mucha, d & d character, in front of an urban background, digital painting, concept art, smooth, sharp focus illustration, artstation hq",
				Sampler:           "Euler a",
				Seed:              950170121,
				Size:              &Size{Width: 960, Height: 960},
				Steps:             20,
			},
		},
		{
			// upscaler latent antialiased
			in: "(cords, antenna:1.1), (chrome:1.3), masterpiece, best quality, [detailed], [intricate], digital painting, portrait of a mechwarrior robot with a (laster gun:1.2), (night:1.5), (mechanical wings, spread wings, banners, glider:1.6), space, (radar:1.3), city, street, cars, skyscrapers, traffic lights, nuclear power plant, computers, science fiction, sci-fi, dieselpunk, WH40K, highres, absurdres, sharp focus, realistic shadows, lithograph by John William Waterhouse and Kyoto animation and Yoshitaka Amano and Frank Frazetta\nNegative prompt: lowres, bad anatomy, bad hands, text, error, missing fingers, extra digit, fewer digits, cropped, worst quality, low quality, normal quality, jpeg artifacts, signature, watermark, username, blurry, artist name, simple background, [nude], [comic panels], [monochrome], [usa], [green background]\nSteps: 20, Sampler: DPM++ 2M Karras, CFG scale: 7, Seed: 2395363541, Size: 640x640, Model hash: 53d4559a, Model: elldrethSLucidMix_v10, Denoising strength: 0, Hires upscale: 2, Hires steps: 1, Hires upscaler: Latent (antialiased)",
			out: Metadata{
				CfgScale:          7,
				DenoisingStrength: 0,
				HiresSteps:        1,
				HiresUpscale:      2,
				HiresUpscaler:     "Latent (antialiased)",
				Model:             "elldrethSLucidMix_v10",
				ModelHash:         "53d4559a",
				NegativePrompt:    "lowres, bad anatomy, bad hands, text, error, missing fingers, extra digit, fewer digits, cropped, worst quality, low quality, normal quality, jpeg artifacts, signature, watermark, username, blurry, artist name, simple background, [nude], [comic panels], [monochrome], [usa], [green background]",
				Prompt:            "(cords, antenna:1.1), (chrome:1.3), masterpiece, best quality, [detailed], [intricate], digital painting, portrait of a mechwarrior robot with a (laster gun:1.2), (night:1.5), (mechanical wings, spread wings, banners, glider:1.6), space, (radar:1.3), city, street, cars, skyscrapers, traffic lights, nuclear power plant, computers, science fiction, sci-fi, dieselpunk, WH40K, highres, absurdres, sharp focus, realistic shadows, lithograph by John William Waterhouse and Kyoto animation and Yoshitaka Amano and Frank Frazetta",
				Sampler:           "DPM++ 2M Karras",
				Seed:              2395363541,
				Size:              &Size{Width: 640, Height: 640},
				Steps:             20,
			},
		},
		{
			// Extra noise = 0.2
			in: "hakurei reimu, (realistic, 3d:0.7), 1girl, portrait, close-up, red eyes, brown hair, hair bow, light smile, closed mouth, white background\nNegative prompt: lowres, bad anatomy, bad hands, text, error, missing fingers, extra digit, fewer digits, cropped, worst quality, low quality, normal quality, jpeg artifacts, signature, watermark, username, blurry\nSteps: 20, Sampler: Euler a, CFG scale: 7, Seed: 903543336, Size: 512x512, Model hash: fdf0096972, VAE hash: c6a580b13a, Denoising strength: 0.45, Clip skip: 2, Hires upscale: 2, Hires steps: 30, Hires upscaler: 4x-UltraSharp, Image noise multiplier: 0.2",
			out: Metadata{
				CfgScale:             7,
				ClipSkip:             2,
				DenoisingStrength:    0.45,
				HiresSteps:           30,
				HiresUpscale:         2,
				HiresUpscaler:        "4x-UltraSharp",
				ImageNoiseMultiplier: 0.2,
				ModelHash:            "fdf0096972",
				NegativePrompt:       "lowres, bad anatomy, bad hands, text, error, missing fingers, extra digit, fewer digits, cropped, worst quality, low quality, normal quality, jpeg artifacts, signature, watermark, username, blurry",
				Prompt:               "hakurei reimu, (realistic, 3d:0.7), 1girl, portrait, close-up, red eyes, brown hair, hair bow, light smile, closed mouth, white background",
				Sampler:              "Euler a",
				Seed:                 903543336,
				Size:                 &Size{Width: 512, Height: 512},
				Steps:                20,
				VaeHash:              "c6a580b13a",
			},
		},
	}

	for i, c := range tc {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			meta, err := ParseParameters(c.in)
			require.NoError(t, err)
			assert.Equal(t, c.out, meta)
		})
	}
}

func TestDecodeMetadata_SD_CPP(t *testing.T) {
	tc := []testCase{
		{
			in: "A sunflower field\nSteps: 30, CFG scale: 5.000000, Guidance: 3.500000, Eta: 0.000000, Seed: 197583933, Size: 512x512, Model: v1-5-pruned-emaonly.safetensors, RNG: cuda, Sampler: dpm++2mv2 karras, Version: stable-diffusion.cpp",
			out: Metadata{
				CfgScale: 5,
				Eta:      0,
				Guidance: 3.5,
				Model:    "v1-5-pruned-emaonly.safetensors",
				Prompt:   "A sunflower field",
				Rng:      "cuda",
				Sampler:  "dpm++2mv2 karras",
				Seed:     197583933,
				Size:     &Size{Width: 512, Height: 512},
				Steps:    30,
				Version:  "stable-diffusion.cpp",
			},
		},
		{
			in: "A sunflower field\nNegative prompt: Blue sky\nSteps: 30, CFG scale: 5.000000, Guidance: 3.500000, Eta: 0.000000, Seed: 1375127038, Size: 512x512, Model: v1-5-pruned-emaonly.safetensors, RNG: cuda, Sampler: dpm++2mv2 karras, Version: stable-diffusion.cpp",
			out: Metadata{
				CfgScale:       5,
				Eta:            0,
				Guidance:       3.5,
				Model:          "v1-5-pruned-emaonly.safetensors",
				NegativePrompt: "Blue sky",
				Prompt:         "A sunflower field",
				Rng:            "cuda",
				Sampler:        "dpm++2mv2 karras",
				Seed:           1375127038,
				Size:           &Size{Width: 512, Height: 512},
				Steps:          30,
				Version:        "stable-diffusion.cpp",
			},
		},
		{
			in: "score_9, score_8_up, score_7_up, sunflower field, (poppy seeds:1.2) \n<lora:Agriculture_V1:1> <lora:SDXL/size_slider_v1:1.7>\nNegative prompt: score_6, score_5, score_4, corn\nSteps: 30, CFG scale: 5.000000, Guidance: 3.500000, Eta: 0.000000, Seed: 1869977377, Size: 1024x1024, Model: sdxl.safetensors, RNG: cuda, Sampler: euler_a karras, VAE: sdxl-vae-fp16-fix.safetensors, Version: stable-diffusion.cpp",
			out: Metadata{
				CfgScale:       5,
				Eta:            0,
				Guidance:       3.5,
				Loras:          Loras{{Name: "Agriculture_V1", Weight: 1}, {Name: "SDXL/size_slider_v1", Weight: 1.7}},
				Model:          "sdxl.safetensors",
				NegativePrompt: "score_6, score_5, score_4, corn",
				Prompt:         "score_9, score_8_up, score_7_up, sunflower field, (poppy seeds:1.2) ,<lora:Agriculture_V1:1> <lora:SDXL/size_slider_v1:1.7>",
				Rng:            "cuda",
				Sampler:        "euler_a karras",
				Seed:           1869977377,
				Size:           &Size{Width: 1024, Height: 1024},
				Steps:          30,
				Vae:            "sdxl-vae-fp16-fix.safetensors",
				Version:        "stable-diffusion.cpp",
			},
		},
		{
			in: "Person in a pirate costume\nSteps: 30, CFG scale: 1.000000, Guidance: 3.500000, Eta: 0.000000, Seed: 1177101575, Size: 512x512, Model: , RNG: cuda, Sampler: euler discrete, TE: clip_l.safetensors, TE: t5-v1_1-xxl-encoder-Q3_K_S.gguf, Unet: PJ0_385_exclusiveTA_00001_BF16_Q4_K_S.gguf, VAE: ae.safetensors, Version: stable-diffusion.cpp",
			out: Metadata{
				CfgScale:    1,
				Guidance:    3.5,
				Model:       "",
				Prompt:      "Person in a pirate costume",
				Rng:         "cuda",
				Sampler:     "euler discrete",
				Seed:        1177101575,
				Size:        &Size{Width: 512, Height: 512},
				Steps:       30,
				TextEncoder: "clip_l.safetensors, t5-v1_1-xxl-encoder-Q3_K_S.gguf",
				Unet:        "PJ0_385_exclusiveTA_00001_BF16_Q4_K_S.gguf",
				Vae:         "ae.safetensors",
				Version:     "stable-diffusion.cpp",
			},
		},
		// stand-alone upscale mode: https://github.com/leejet/stable-diffusion.cpp/pull/865
		{
			in: "\nSteps: 20, CFG scale: 7.000000, Guidance: 3.500000, Eta: 0.000000, Seed: 42, Size: 338x338, Model: , RNG: cuda, Sampler: default, Version: stable-diffusion.cpp",
			out: Metadata{
				CfgScale: 7,
				Guidance: 3.5,
				Model:    "",
				Prompt:   "",
				Rng:      "cuda",
				Sampler:  "default",
				Seed:     42,
				Size:     &Size{Width: 338, Height: 338},
				Steps:    20,
				Version:  "stable-diffusion.cpp",
			},
		},
	}

	for i, c := range tc {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			meta, err := ParseParameters(c.in)
			require.NoError(t, err)
			assert.Equal(t, c.out, meta)
		})
	}
}

func TestDecodeMetadata_Fooocus(t *testing.T) {
	tc := []testCase{
		{
			in: "cinematic still A sunflower field . emotional, harmonious, vignette, 4k epic detailed, shot on kodak, 35mm photo, sharp focus, high budget, cinemascope, moody, epic, gorgeous, film grain, grainy, A sunflower field, highly detailed, magic, peaceful, flowing, beautiful, atmosphere, radiant, magical, sharp focus, very coherent, intricate, elegant, epic, colorful, amazing composition, cinematic, artistic, fine detail, professional, clear, joyful, unique, expressive, cute, iconic, best, vivid, awesome, perfect, ambient background, pristine, creative\nNegative prompt: (worst quality, low quality, normal quality, lowres, low details, oversaturated, undersaturated, overexposed, underexposed, grayscale, bw, bad photo, bad photography, bad art:1.4), (watermark, signature, text font, username, error, logo, words, letters, digits, autograph, trademark, name:1.2), (blur, blurry, grainy), morbid, ugly, asymmetrical, mutated malformed, mutilated, poorly lit, bad shadow, draft, cropped, out of frame, cut off, censored, jpeg artifacts, out of focus, glitch, duplicate, (airbrushed, cartoon, anime, semi-realistic, cgi, render, blender, digital art, manga, amateur:1.3), (3D ,3D Game, 3D Game Scene, 3D Character:1.1), (bad hands, bad anatomy, bad body, bad face, bad teeth, bad arms, bad legs, deformities:1.3), anime, cartoon, graphic, (blur, blurry, bokeh), text, painting, crayon, graphite, abstract, glitch, deformed, mutated, ugly, disfigured\nSteps: 30, Sampler: DPM++ 2M SDE Karras, Seed: 127589946317439009, Size: 512x512, CFG scale: 4, Sharpness: 2, ADM Guidance: \"(1.5, 0.8, 0.3)\", Model: juggernautXL_v8Rundiffusion, Model hash: aeb7e9e689, Performance: Speed, Scheduler: karras, VAE: Default (model), Raw prompt: A sunflower field, Raw negative prompt: , Clip skip: 2, Lora hashes: \"sd_xl_offset_example-lora_1.0: 4852686128\", Lora weights: \"sd_xl_offset_example-lora_1.0: 0.1\", Version: Fooocus v2.5.5",
			out: Metadata{
				CfgScale:       4,
				ClipSkip:       2,
				Model:          "juggernautXL_v8Rundiffusion",
				ModelHash:      "aeb7e9e689",
				NegativePrompt: "(worst quality, low quality, normal quality, lowres, low details, oversaturated, undersaturated, overexposed, underexposed, grayscale, bw, bad photo, bad photography, bad art:1.4), (watermark, signature, text font, username, error, logo, words, letters, digits, autograph, trademark, name:1.2), (blur, blurry, grainy), morbid, ugly, asymmetrical, mutated malformed, mutilated, poorly lit, bad shadow, draft, cropped, out of frame, cut off, censored, jpeg artifacts, out of focus, glitch, duplicate, (airbrushed, cartoon, anime, semi-realistic, cgi, render, blender, digital art, manga, amateur:1.3), (3D ,3D Game, 3D Game Scene, 3D Character:1.1), (bad hands, bad anatomy, bad body, bad face, bad teeth, bad arms, bad legs, deformities:1.3), anime, cartoon, graphic, (blur, blurry, bokeh), text, painting, crayon, graphite, abstract, glitch, deformed, mutated, ugly, disfigured",
				Prompt:         "cinematic still A sunflower field . emotional, harmonious, vignette, 4k epic detailed, shot on kodak, 35mm photo, sharp focus, high budget, cinemascope, moody, epic, gorgeous, film grain, grainy, A sunflower field, highly detailed, magic, peaceful, flowing, beautiful, atmosphere, radiant, magical, sharp focus, very coherent, intricate, elegant, epic, colorful, amazing composition, cinematic, artistic, fine detail, professional, clear, joyful, unique, expressive, cute, iconic, best, vivid, awesome, perfect, ambient background, pristine, creative",
				Sampler:        "DPM++ 2M SDE Karras",
				Seed:           127589946317439009,
				Size:           &Size{Width: 512, Height: 512},
				Steps:          30,
				Vae:            "Default (model)",
				Version:        "Fooocus v2.5.5",
				// Unsupported fields:
				// ADM Guidance: "(1.5, 0.8, 0.3)",
				// Performance: Speed,
				// Scheduler: karras,
				// Raw prompt: A sunflower field,
				// Raw negative prompt: ,
				// Loras:    Loras{{Name: "sd_xl_offset_example-lora_1.0", Weight: 0.1}},
				// Lora hashes: "sd_xl_offset_example-lora_1.0: 4852686128",
				// Lora weights: "sd_xl_offset_example-lora_1.0: 0.1",
				// Sharpness: 2,
			},
		},
	}

	for i, c := range tc {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			meta, err := ParseParameters(c.in)
			require.NoError(t, err)
			assert.Equal(t, c.out, meta)
		})
	}
}

func TestDecodeMetadata_Fooocus_JSON_Error(t *testing.T) {
	tc := []string{
		`{"adm_guidance": "(1.5, 0.8, 0.3)", "base_model": "juggernautXL_v8Rundiffusion", "base_model_hash": "aeb7e9e689", "clip_skip": 2}`,
	}

	for i, c := range tc {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			_, err := ParseParameters(c)
			assert.Error(t, err)
		})
	}
}
