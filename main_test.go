package metadata

import (
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/fkleon/fooocus-metadata/fooocus"
	_ "github.com/fkleon/fooocus-metadata/fooocusplus"
	_ "github.com/fkleon/fooocus-metadata/ruinedfooocus"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Configure logging during testing
func TestMain(m *testing.M) {
	slog.SetLogLoggerLevel(slog.LevelWarn)
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestExtractOne_Fooocus(t *testing.T) {
	const testpath = "./fooocus/testdata/"
	var files = []string{
		"fooocus-meta.png",
	}

	for _, file := range files {
		t.Run(path.Base(file), func(t *testing.T) {
			path := filepath.Join(testpath, file)
			meta, err := ExtractFromFile(path)
			require.NoError(t, err)
			assert.Equal(t, "Fooocus", meta.Source)
			assert.Equal(t, "Fooocus v2.5.5", meta.Params.Version())
			assert.Equal(t, "juggernautXL_v8Rundiffusion", meta.Params.Model())
			assert.NotZero(t, meta.Params.Raw())
		})
	}
}

func TestExtractOne_FooocusPlus(t *testing.T) {
	const testpath = "./fooocusplus/testdata/"
	var files = []string{
		"fooocusplus-meta.png",
	}

	for _, file := range files {
		t.Run(path.Base(file), func(t *testing.T) {
			path := filepath.Join(testpath, file)
			meta, err := ExtractFromFile(path)
			require.NoError(t, err)
			assert.Equal(t, "FooocusPlus", meta.Source)
			assert.Equal(t, "FooocusPlus 1.0.0", meta.Params.Version())
			assert.Equal(t, "elsewhereXL_v10", meta.Params.Model())
			assert.NotZero(t, meta.Params.Raw())
		})
	}
}

func TestExtractOne_RuinedFooocus(t *testing.T) {
	const testpath = "./ruinedfooocus/testdata/"
	var files = []string{
		"ruinedfooocus-meta.png",
	}

	for _, file := range files {
		t.Run(path.Base(file), func(t *testing.T) {
			path := filepath.Join(testpath, file)
			meta, err := ExtractFromFile(path)
			require.NoError(t, err)
			assert.Equal(t, "RuinedFooocus", meta.Source)
			assert.Equal(t, "RuinedFooocus", meta.Params.Version())
			assert.Equal(t, "sd_xl_base_1.0_0.9vae", meta.Params.Model())
			assert.NotZero(t, meta.Params.Raw())
		})
	}
}

func TestExtractMetadata_Fooocus(t *testing.T) {
	const testpath = "./fooocus/testdata/"
	testCases := []struct {
		file     string
		source   string
		software string
	}{
		{"fooocus-meta.png", "Fooocus", "Fooocus v2.5.5"},
		{"fooocus-meta.jpeg", "Fooocus", "Fooocus v2.5.5"},
		{"fooocus-meta.webp", "Fooocus", "Fooocus v2.5.5"},
	}

	for _, tc := range testCases {
		t.Run(tc.file, func(t *testing.T) {
			path := filepath.Join(testpath, tc.file)
			meta, err := ExtractFromFile(path)

			require.NoError(t, err)
			require.NotNil(t, meta)

			assert.Equal(t, tc.source, meta.Source)
			assert.Equal(t, tc.software, meta.Params.Version())
			assert.Zero(t, meta.Created)
		})
	}
}

func TestExtractMetadata_FooocusPlus(t *testing.T) {
	const testpath = "./fooocusplus/testdata/"
	testCases := []struct {
		file     string
		source   string
		software string
	}{
		{"fooocusplus-meta.png", "FooocusPlus", "FooocusPlus 1.0.0"},
		{"fooocusplus-meta.jpg", "FooocusPlus", "FooocusPlus 1.0.0"},
		{"fooocusplus-meta.webp", "FooocusPlus", "FooocusPlus 1.0.0"},
	}

	for _, tc := range testCases {
		t.Run(tc.file, func(t *testing.T) {
			path := filepath.Join(testpath, tc.file)
			meta, err := ExtractFromFile(path)

			require.NoError(t, err)
			require.NotNil(t, meta)

			assert.Equal(t, tc.source, meta.Source)
			assert.Equal(t, tc.software, meta.Params.Version())
			assert.Zero(t, meta.Created)
		})
	}
}

func TestExtractMetadata_RuinedFooocus(t *testing.T) {
	const testpath = "./ruinedfooocus/testdata/"
	testCases := []struct {
		file     string
		source   string
		software string
	}{
		{"ruinedfooocus-meta.png", "RuinedFooocus", "RuinedFooocus"},
	}

	for _, tc := range testCases {
		t.Run(tc.file, func(t *testing.T) {
			path := filepath.Join(testpath, tc.file)
			meta, err := ExtractFromFile(path)

			require.NoError(t, err)
			require.NotNil(t, meta)

			assert.Equal(t, tc.source, meta.Source)
			assert.Equal(t, tc.software, meta.Params.Version())
			assert.Zero(t, meta.Created)
		})
	}
}

func TestExtractCreatedTime(t *testing.T) {
	filenamePattern := "2024-01-05_23-11-48_9167_*.png"
	expectedCreatedTime := time.Date(2024, time.January, 5, 23, 11, 48, 0, time.UTC)

	files, err := filepath.Glob("./*/testdata/*fooocus*.png")
	require.NoError(t, err)

	for _, tc := range files {
		t.Run(path.Base(tc), func(t *testing.T) {
			in, err := os.Open(tc)
			out := createTemp(t, filenamePattern)

			_, err = io.Copy(out, in)
			require.NoError(t, err)

			meta, err := ExtractFromFile(out.Name())
			require.NoError(t, err)
			require.NotNil(t, meta)
			assert.Equal(t, expectedCreatedTime, meta.Created)
		})
	}
}

/*
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
*/

// Create a temp file and register a callback to clean it up after the test run
func createTemp(t *testing.T, pattern string) *os.File {
	target, err := os.CreateTemp("", pattern)
	require.NoError(t, err)
	t.Cleanup(func() {
		os.Remove(target.Name())
	})
	return target
}
