package metadata

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/fkleon/fooocus-metadata/internal/fooocus"
	"github.com/fkleon/fooocus-metadata/internal/fooocusplus"
	"github.com/fkleon/fooocus-metadata/internal/ruinedfooocus"
	"github.com/stretchr/testify/require"
)

// Configure logging during testing
func TestMain(m *testing.M) {
	slog.SetLogLoggerLevel(slog.LevelWarn)
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestExtractOne_Fooocus(t *testing.T) {
	const testpath = "./internal/fooocus/testdata/"
	var files = []string{
		"fooocus-meta.png",
	}

	for _, file := range files {
		t.Run(path.Base(file), func(t *testing.T) {
			path := filepath.Join(testpath, file)
			meta, err := ExtractOne[fooocus.Metadata](path)
			require.NoError(t, err)
			require.IsType(t, fooocus.Metadata{}, meta)
		})
	}
}

func TestExtractOne_FooocusPlus(t *testing.T) {
	const testpath = "./internal/fooocusplus/testdata/"
	var files = []string{
		"fooocusplus-meta.png",
	}

	for _, file := range files {
		t.Run(path.Base(file), func(t *testing.T) {
			path := filepath.Join(testpath, file)
			meta, err := ExtractOne[fooocusplus.Metadata](path)
			require.NoError(t, err)
			require.IsType(t, fooocusplus.Metadata{}, meta)
		})
	}
}

func TestExtractOne_RuinedFooocus(t *testing.T) {
	const testpath = "./internal/ruinedfooocus/testdata/"
	var files = []string{
		"ruinedfooocus-meta.png",
	}

	for _, file := range files {
		t.Run(path.Base(file), func(t *testing.T) {
			path := filepath.Join(testpath, file)
			meta, err := ExtractOne[ruinedfooocus.Metadata](path)
			require.NoError(t, err)
			require.IsType(t, ruinedfooocus.Metadata{}, meta)
		})
	}
}

func TestExtractMetadata_Fooocus(t *testing.T) {
	const testpath = "./internal/fooocus/testdata/"
	testCases := []struct {
		file             string
		expectedType     interface{}
		expectedSoftware string
	}{
		{"fooocus-meta.png", fooocus.Metadata{}, "Fooocus v2.5.5"},
		{"fooocus-meta.jpeg", fooocus.Metadata{}, "Fooocus v2.5.5"},
		{"fooocus-meta.webp", fooocus.Metadata{}, "Fooocus v2.5.5"},
	}

	for _, tc := range testCases {
		t.Run(tc.file, func(t *testing.T) {
			path := filepath.Join(testpath, tc.file)
			params, err := ExtractFromFile(path)

			require.NoError(t, err)
			require.NotNil(t, params)

			require.IsType(t, tc.expectedType, params.Raw())
			require.Equal(t, tc.expectedSoftware, params.Software())
			require.True(t, params.CreatedTime().IsZero())
		})
	}
}

func TestExtractMetadata_FooocusPlus(t *testing.T) {
	const testpath = "./internal/fooocusplus/testdata/"
	testCases := []struct {
		file             string
		expectedType     interface{}
		expectedSoftware string
	}{
		{"fooocusplus-meta.png", fooocusplus.Metadata{}, "FooocusPlus 1.0.0"},
		{"fooocusplus-meta.jpg", fooocusplus.Metadata{}, "FooocusPlus 1.0.0"},
		{"fooocusplus-meta.webp", fooocusplus.Metadata{}, "FooocusPlus 1.0.0"},
	}

	for _, tc := range testCases {
		t.Run(tc.file, func(t *testing.T) {
			path := filepath.Join(testpath, tc.file)
			params, err := ExtractFromFile(path)

			require.NoError(t, err)
			require.NotNil(t, params)

			require.IsType(t, tc.expectedType, params.Raw())
			require.Equal(t, tc.expectedSoftware, params.Software())
			require.True(t, params.CreatedTime().IsZero())
		})
	}
}

func TestExtractMetadata_RuinedFooocus(t *testing.T) {
	const testpath = "./internal/ruinedfooocus/testdata/"
	testCases := []struct {
		file             string
		expectedType     interface{}
		expectedSoftware string
	}{
		{"ruinedfooocus-meta.png", ruinedfooocus.Metadata{}, "RuinedFooocus"},
	}

	for _, tc := range testCases {
		t.Run(tc.file, func(t *testing.T) {
			path := filepath.Join(testpath, tc.file)
			params, err := ExtractFromFile(path)

			require.NoError(t, err)
			require.NotNil(t, params)

			require.IsType(t, tc.expectedType, params.Raw())
			require.Equal(t, tc.expectedSoftware, params.Software())
			require.True(t, params.CreatedTime().IsZero())
		})
	}
}

func TestExtractCreatedTime(t *testing.T) {
	filenamePattern := "2024-01-05_23-11-48_9167_*.png"
	expectedCreatedTime := time.Date(2024, time.January, 5, 23, 11, 48, 0, time.UTC)

	files, err := filepath.Glob("./internal/*/testdata/*fooocus*.png")
	require.NoError(t, err)

	for _, tc := range files {
		t.Run(path.Base(tc), func(t *testing.T) {
			in, err := os.Open(tc)
			out := createTemp(t, filenamePattern)

			_, err = io.Copy(out, in)
			require.NoError(t, err)

			params, err := ExtractFromFile(out.Name())
			require.NoError(t, err)
			require.NotNil(t, params)
			require.Equal(t, expectedCreatedTime, params.CreatedTime())
		})
	}
}

func TestEmbedMetadata_Fooocus(t *testing.T) {

	target := createTemp(t, "out.fooocus.*.png")

	err := EmbedIntoFile(target.Name(), fooocus.Metadata{
		Version: "Fooocus Metadata",
	})
	require.NoError(t, err)
}

func TestEmbedMetadata_FooocusPlus(t *testing.T) {

	target := createTemp(t, "out.fooocusplus.*.png")

	err := EmbedIntoFile(target.Name(), fooocusplus.Metadata{
		Version: "FooocusPlus Metadata",
	})
	require.NoError(t, err)
}

func TestEmbedMetadata_RuinedFooocus(t *testing.T) {

	target := createTemp(t, "out.ruinedfooocus.*.png")

	err := EmbedIntoFile(target.Name(), ruinedfooocus.Metadata{
		Version: "RuinedFooocus Metadata",
	})
	require.NoError(t, err)
}

func TestEmbedMetadata_WithSource(t *testing.T) {
	const testpath = "./testdata/"
	testCases := []string{
		"sample.png",
		"sample.jpg",
		"sample.webp",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			path := filepath.Join(testpath, tc)
			source, err := os.Open(path)
			require.NoError(t, err)

			target := createTemp(t, fmt.Sprintf("out.%s.*.png", tc))

			err = Embed(EmbedOptions{
				Source: source,
				Target: target,
			}, fooocus.Metadata{
				Version: "Fooocus Metadata",
			})
			require.NoError(t, err)
		})
	}
}

// Create a temp file and register a callback to clean it up after the test run
func createTemp(t *testing.T, pattern string) *os.File {
	target, err := os.CreateTemp("", pattern)
	require.NoError(t, err)
	t.Cleanup(func() {
		os.Remove(target.Name())
	})
	return target
}
