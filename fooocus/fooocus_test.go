package fooocus

import (
	"bytes"
	"image"
	"image/png"
	"strconv"
	"testing"

	"github.com/bep/imagemeta"
	"github.com/fkleon/fooocus-metadata/types"
	pngembed "github.com/sabhiram/png-embed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractMetadataFromPNG(t *testing.T) {

	var pngData map[string]imagemeta.TagInfo = make(map[string]imagemeta.TagInfo)

	pngData["fooocus_scheme"] =
		imagemeta.TagInfo{
			Source:    0,
			Namespace: "PNG/tEXt",
			Tag:       "fooocus_scheme",
			Value:     Fooocus.String(),
		}
	pngData["parameters"] = imagemeta.TagInfo{
		Source:    0,
		Namespace: "PNG/tEXt",
		Tag:       "parameters",
		Value:     metaV23Json,
	}

	extractor := NewFooocusMetadataExtractor()
	fooocusData, err := extractor.Decode(types.ImageMetadataContext{
		EmbeddedMetadata: pngData,
	})
	require.NoError(t, err)
	require.Equal(t, *metaV23, fooocusData)
}

func TestExtractMetadataFromExif(t *testing.T) {

	var exifData imagemeta.Tags

	exifData.Add(imagemeta.TagInfo{
		Source: imagemeta.EXIF,
		Tag:    "Software",
		Value:  "Fooocus v2.5.5",
	})
	exifData.Add(imagemeta.TagInfo{
		Source: imagemeta.EXIF,
		Tag:    "MakerNoteApple",
		Value:  Fooocus.String(),
	})
	exifData.Add(imagemeta.TagInfo{
		Source: imagemeta.EXIF,
		Tag:    "UserComment",
		Value:  metaV23Json,
	})

	extractor := NewFooocusMetadataExtractor()
	fooocusData, err := extractor.Decode(types.ImageMetadataContext{
		EmbeddedMetadata: exifData.EXIF(),
	})
	require.NoError(t, err)
	require.Equal(t, *metaV23, fooocusData)
}

func TestExtractMetadataFromSidecar(t *testing.T) {
	extractor := NewFooocusMetadataExtractor()

	ctx := types.ImageMetadataContext{
		Filepath: "./testdata/fooocus-meta.png",
	}
	structMeta, err := extractor.Extract(ctx)
	require.NoError(t, err)
	require.NotZero(t, structMeta.Params.Raw())
	require.Equal(t, "juggernautXL_v8Rundiffusion", structMeta.Params.Model())
}

func TestEmbedMetadataIntoPNG_CopyWrite(t *testing.T) {
	writer := NewFooocusMetadataWriter()

	// Source image is a 1x1 pixel PNG
	source := &bytes.Buffer{}
	err := png.Encode(source, image.NewRGBA(image.Rect(0, 0, 1, 1)))
	require.NoError(t, err)

	target := &bytes.Buffer{}

	err = writer.CopyWrite(source, target, *metaV23)
	require.NoError(t, err)

	// Expect resulting PNG to have 2 embedded chunks
	data, err := pngembed.Extract(target.Bytes())
	require.NoError(t, err)
	require.Equal(t, 2, len(data))
	require.Equal(t, []byte(Fooocus.String()), data["fooocus_scheme"])
	require.Contains(t, data, "parameters")

	// Expect resulting PNG to have same dimensions as source
	info, _, err := image.Decode(target)
	require.NoError(t, err)
	require.Equal(t, image.Pt(1, 1), info.Bounds().Max)
}

func TestEmbedMetadataIntoPNG_Write(t *testing.T) {
	writer := NewFooocusMetadataWriter()
	target := &bytes.Buffer{}

	// Source image is the default template
	err := writer.Write(target, *metaV23)
	require.NoError(t, err)

	// Expect resulting PNG to have 5 embedded chunks:
	// 3 from template, 2 from metadata
	data, err := pngembed.Extract(target.Bytes())
	require.NoError(t, err)
	require.Equal(t, 5, len(data))
	require.Equal(t, []byte(Fooocus.String()), data["fooocus_scheme"])
	require.Contains(t, data, "parameters")

	// Expect resulting PNG to have same dimensions as template
	info, _, err := image.Decode(target)
	require.NoError(t, err)
	require.Equal(t, image.Pt(240, 85), info.Bounds().Max)
}

func TestAdapter(t *testing.T) {
	testCases := []struct {
		meta  Metadata
		model string
	}{
		{*metaV23, "juggernautXL_v8Rundiffusion"},
		{*metaV23Alt, "ponyDiffusionV6XL_v6TurboDPOMerge"},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			param := Parameters{
				Metadata: tc.meta,
			}
			assert.Equal(t, tc.model, param.Model())
		})
	}
}
