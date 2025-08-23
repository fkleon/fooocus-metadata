package types

import (
	"bytes"
	"image"
	"os"
	"testing"

	_ "image/jpeg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmbedWithoutSource(t *testing.T) {
	// Test embedding without a source file
	// This should write a copy of the template PNG
	// with the metadata embedded
	writer := NewPngMetadataWriter()

	values := make(map[string]interface{})
	values["fooocus_scheme"] = "Fooocus"
	values["parameters"] = nil

	var buf bytes.Buffer

	err := writer.Embed(nil, &buf, values)
	assert.NoError(t, err)
	assert.NotEmpty(t, buf)

	image, format, err := image.Decode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, "png", format, "Expected PNG format after embedding metadata")
	// Expect dimensions of template PNG
	assert.Equal(t, image.Bounds().Dx(), 240)
	assert.Equal(t, image.Bounds().Dy(), 85)
}

func TestEmbedWithSource(t *testing.T) {
	// Test embedding with a source file
	// This should write a copy of the source (converted to PNG)
	// with the metadata embedded

	writer := NewPngMetadataWriter()

	values := make(map[string]interface{})
	values["fooocus_scheme"] = "Fooocus"
	values["parameters"] = nil

	source, err := os.Open("../fooocus/testdata/fooocus-meta.jpeg")
	require.NoError(t, err)
	defer source.Close()

	var buf bytes.Buffer

	err = writer.Embed(source, &buf, values)
	assert.NoError(t, err)
	assert.NotEmpty(t, buf)

	image, format, err := image.Decode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, "png", format, "Expected PNG format after embedding metadata")
	// Expect dimensions of source image
	assert.Equal(t, image.Bounds().Dx(), 512)
	assert.Equal(t, image.Bounds().Dy(), 512)
}
