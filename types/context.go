package types

import "github.com/bep/imagemeta"

type ImageMetadataContext struct {

	// Original filepath of the image, used to locate any sidecar files
	// if supported.
	// This is optional. If the image is piped or memory-only, extractors
	// should gracefully skip sidecar logic.
	Filepath string

	// MIME type of the image, e.g. "image/jpeg", "image/png"
	MIME string

	// All embedded metadata extracted from the image, usually
	// from EXIF blocks or PNG tEXt chunks.
	EmbeddedMetadata map[string]imagemeta.TagInfo
}
