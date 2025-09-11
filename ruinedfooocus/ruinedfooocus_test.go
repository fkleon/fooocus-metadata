package ruinedfooocus

import (
	"testing"

	"github.com/bep/imagemeta"
	"github.com/fkleon/fooocus-metadata/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractMetadataFromPNG(t *testing.T) {

	var pngData = make(map[string]imagemeta.TagInfo)

	pngData["parameters"] = imagemeta.TagInfo{
		Source:    0,
		Namespace: "PNG/tEXt",
		Tag:       "parameters",
		Value:     metaJson,
	}

	extractor := NewRuinedFooocusMetadataExtractor()
	fooocusData, err := extractor.Decode(types.ImageMetadataContext{
		EmbeddedMetadata: pngData,
	})
	require.NoError(t, err)
	require.Equal(t, *meta, fooocusData)
}

func TestAdapter(t *testing.T) {
	// With file extension
	param := Parameters{
		Metadata: *meta,
	}
	assert.Equal(t, "sd_xl_base_1.0_0.9vae", param.Model())

	// With path
	param.BaseModel = "Pony/ponysdxl.safetensors"
	assert.Equal(t, "ponysdxl", param.Model())
}
