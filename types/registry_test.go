package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeWithoutReader(t *testing.T) {
	ctx := ImageMetadataContext{
		Filepath: "testdata/sample.jpg",
		MIME:     "image/jpeg",
	}
	_, err := Decode(ctx)
	require.Error(t, err, "no metadata found")
}

func TestDecodeWithReader(t *testing.T) {
	ctx := ImageMetadataContext{
		Filepath: "testdata/sample.jpg",
		MIME:     "image/jpeg",
	}

	source := "TestSource"
	reader := func(ctx ImageMetadataContext) (StructuredMetadata, error) {
		return StructuredMetadata{
			Source: source,
		}, nil
	}

	RegisterReader(source, reader)

	meta, err := Decode(ctx)
	require.NoError(t, err)
	require.Equal(t, source, meta.Source)
}

func TestDecodeWithMultipleReaders(t *testing.T) {

	RegisterReader("TestErrorSource", func(ctx ImageMetadataContext) (StructuredMetadata, error) {
		return StructuredMetadata{}, fmt.Errorf("an error occurred")
	})
	RegisterReader("TestSource", func(ctx ImageMetadataContext) (StructuredMetadata, error) {
		return StructuredMetadata{
			Source: "TestSource",
		}, nil
	})

	meta, err := Decode(ImageMetadataContext{})
	require.NoError(t, err)
	require.Equal(t, "TestSource", meta.Source)
}
