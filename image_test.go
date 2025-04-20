package fooocus

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Disable logging during testing
func TestMain(m *testing.M) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestExtractExif(t *testing.T) {
	// Original released under Public Domain: https://commons.wikimedia.org/wiki/File:Steveston_dusk.JPG
	file, err := os.Open("testdata/sample.jpg")
	require.NoError(t, err)

	exifData, err := extractExif(file, "image/jpeg")
	require.NoError(t, err)
	assert.NotNil(t, exifData)

	exifVersion, _ := exifData.EXIF()["ExifVersion"]
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
	image, err := NewImageInfo(path)
	require.NoError(t, err)

	assert.Equal(t, "image/jpeg", image.MIME)
	assert.NotNil(t, image.exif)
	assert.Nil(t, image.pngText)
	assert.Nil(t, image.FooocusMetadata)
}

func TestExtractImageInfo_PNG(t *testing.T) {
	path := "testdata/sample.png"
	image, err := NewImageInfo(path)
	require.NoError(t, err)

	assert.Equal(t, "image/png", image.MIME)
	assert.Nil(t, image.exif)
	assert.NotNil(t, image.pngText)
	assert.Nil(t, image.FooocusMetadata)
}

func TestExtractImageInfo_WEBP(t *testing.T) {
	path := "testdata/sample.webp"
	image, err := NewImageInfo(path)
	require.NoError(t, err)

	assert.Equal(t, "image/webp", image.MIME)
	assert.NotNil(t, image.exif)
	assert.Nil(t, image.pngText)
	assert.Nil(t, image.FooocusMetadata)
}

func TestExtractFooocusMetadata(t *testing.T) {
	testCases := []struct {
		file    string
		hasExif bool
	}{
		{"fooocus-meta.png", false},
		{"fooocus-meta.jpeg", true},
		{"fooocus-meta.webp", true},
	}

	for _, tc := range testCases {
		t.Run(tc.file, func(t *testing.T) {
			path := filepath.Join("testdata", tc.file)
			image, err := NewImageInfo(path)
			require.NoError(t, err)

			if tc.hasExif {
				assert.NotNil(t, image.exif)
				assert.Nil(t, image.pngText)
			} else {
				assert.Nil(t, image.exif)
				assert.NotNil(t, image.pngText)
			}

			require.NotNil(t, image.FooocusMetadata)
			assert.Equal(t, "Fooocus v2.5.5", image.FooocusMetadata.Version)
		})
	}
}

func TestEmbedFooocusMetadataAsPng_WithSource(t *testing.T) {
	source, err := os.Open("testdata/sample.png")
	assert.NoError(t, err)

	target, err := os.CreateTemp("", "out.source.*.png")
	assert.NoError(t, err)

	err = EmbedMetadataAsPngText(source, target, meta)
	assert.NoError(t, err)
}

func TestEmbedFooocusMetadataAsPng_WithoutSource(t *testing.T) {
	target, err := os.CreateTemp("", "out.template.*.png")
	assert.NoError(t, err)

	err = EmbedMetadataAsPngText(nil, target, meta)
	assert.NoError(t, err)
}
