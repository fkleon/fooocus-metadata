package types

import "github.com/bep/imagemeta"

type ImageMetadataContext struct {

	// Original filepath of the image, used to locate any sidecar files
	// if supported. If the image was read from a stream, this may be empty.
	//
	// Note: This is the original path of the image, not the path of any
	// sidecar file. Sidecar files must be located by the extractor if
	// supported.
	Filepath string

	// MIME type of the image, e.g. "image/jpeg", "image/png"
	MIME string

	// All embedded metadata extracted from the image, usually
	// from EXIF blocks or PNG tEXt chunks.
	EmbeddedMetadata map[string]imagemeta.TagInfo
}
