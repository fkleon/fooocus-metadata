package ruinedfooocus

import (
	"testing"

	"github.com/bep/imagemeta"
	"github.com/fkleon/fooocus-metadata/types"
	"github.com/stretchr/testify/require"
)

func TestExtractMetadataFromPNG(t *testing.T) {

	var pngData map[string]imagemeta.TagInfo = make(map[string]imagemeta.TagInfo)

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
