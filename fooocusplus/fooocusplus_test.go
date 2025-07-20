package fooocusplus

import (
	"testing"

	"github.com/bep/imagemeta"
	"github.com/fkleon/fooocus-metadata/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractMetadataFromPNG(t *testing.T) {
	var pngData map[string]imagemeta.TagInfo = make(map[string]imagemeta.TagInfo)

	pngData["Comment"] = imagemeta.TagInfo{
		Source:    0,
		Namespace: "PNG/tEXt",
		Tag:       "Comment",
		Value:     metaJson,
	}

	extractor := NewFooocusPlusMetadataExtractor()
	fooocusData, err := extractor.Decode(types.ImageMetadataContext{
		EmbeddedMetadata: pngData,
	})
	require.NoError(t, err)
	require.Equal(t, *meta, fooocusData)
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

	extractor := NewFooocusPlusMetadataExtractor()
	fooocusData, err := extractor.Decode(types.ImageMetadataContext{
		EmbeddedMetadata: exifData.EXIF(),
	})
	require.NoError(t, err)
	require.Equal(t, *meta, fooocusData)
}

func TestExtractMetadataFromSidecar(t *testing.T) {
	extractor := NewFooocusPlusMetadataExtractor()

	ctx := types.ImageMetadataContext{
		Filepath: "./testdata/fooocusplus-meta.png",
	}
	structMeta, err := extractor.Extract(ctx)
	require.NoError(t, err)
	require.NotZero(t, structMeta.Params.Raw())
	require.Equal(t, "elsewhereXL_v10", structMeta.Params.Model())
}

func TestAdapter(t *testing.T) {
	param := Parameters{
		Metadata: *meta,
	}
	assert.Equal(t, "elsewhereXL_v10", param.Model())
}
