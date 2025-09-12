package image

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bep/imagemeta"
)

func TestExtractExif(t *testing.T) {
	// Original released under Public Domain: https://commons.wikimedia.org/wiki/File:Steveston_dusk.JPG
	file, err := os.Open("testdata/sample.jpg")
	require.NoError(t, err)

	exifData, err := extractExif(file, "image/jpeg")
	require.NoError(t, err)
	assert.NotNil(t, exifData)

	exifVersion := exifData.EXIF()["ExifVersion"]
	exifVersionStr := exifVersion.Value.(string)
	assert.Equal(t, "0220", exifVersionStr)
}

func TestExtractPNGTextChunks(t *testing.T) {
	file, err := os.Open("testdata/sample.png")
	require.NoError(t, err)

	meta, err := extractPngTextChunks(file)
	require.NoError(t, err)

	assert.Equal(t, map[string]string{
		"date:create":    "2025-04-11T09:41:46+00:00",
		"date:modify":    "2025-04-11T09:41:46+00:00",
		"date:timestamp": "2025-04-11T11:53:39+00:00",
		"Software":       "ImageMaker2000(TM)",
	}, meta)
}

func TestExtractImageInfo_JPEG(t *testing.T) {
	path := "testdata/sample.jpg"
	image, err := NewContextFromFile(path)
	require.NoError(t, err)

	assert.Equal(t, "image/jpeg", image.MIME)
	for _, v := range image.EmbeddedMetadata {
		assert.Equal(t, imagemeta.EXIF, v.Source)
		assert.Contains(t, v.Namespace, "IFD0")
	}
}

func TestExtractImageInfo_PNG(t *testing.T) {
	path := "testdata/sample.png"
	image, err := NewContextFromFile(path)
	require.NoError(t, err)

	assert.Equal(t, "image/png", image.MIME)
	for _, v := range image.EmbeddedMetadata {
		assert.Equal(t, imagemeta.Source(0x0), v.Source)
		assert.Equal(t, "PNG/tEXt", v.Namespace)
	}
}

func TestExtractImageInfo_WEBP(t *testing.T) {
	path := "testdata/sample.webp"
	image, err := NewContextFromFile(path)
	require.NoError(t, err)

	assert.Equal(t, "image/webp", image.MIME)
	for _, v := range image.EmbeddedMetadata {
		assert.Equal(t, imagemeta.EXIF, v.Source)
		assert.Contains(t, v.Namespace, "IFD0")
	}
}
