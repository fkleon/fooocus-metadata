package types

import (
	"fmt"
	"log/slog"
	"sync"
)

type format struct {
	name   string
	decode func(ImageMetadataContext) (StructuredMetadata, error)
	//encode func(GenerationParameters) (error)
}

var (
	formatsMu sync.Mutex
	formats   []format = make([]format, 0, 3)
)

func RegisterReader(name string, decode func(ImageMetadataContext) (StructuredMetadata, error)) {
	formatsMu.Lock()
	formats = append(formats, format{name, decode})
	formatsMu.Unlock()
}

func Decode(ctx ImageMetadataContext) (StructuredMetadata, error) {
	slog.Debug("Decoding metadata", "mime", ctx.MIME, "count", len(ctx.EmbeddedMetadata))

	for _, format := range formats {
		slog.Debug("Trying to decode with", "software", format.name)
		if params, err := format.decode(ctx); err == nil {
			slog.Debug("Found metadata", "software", format.name)
			return params, nil
		}
	}

	return StructuredMetadata{}, fmt.Errorf("no metadata found")
}
