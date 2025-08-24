package types

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormaliseModelName(t *testing.T) {

	testCases := []struct {
		input    string
		expected string
	}{
		{"model.safetensors", "model"},
		{"/path/to/model.safetensors", "model"},
		{"another_model.gguf", "another_model"},
		{"models/model.pt", "model"},
		{"models/model.pth", "model"},
		{"relative/path/model.onnx", "model"},
		{"simplemodel", "simplemodel"},
		{"model.txt", "model.txt"},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			result := NormaliseModelName(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}
