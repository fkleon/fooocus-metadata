package fooocus

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/antchfx/htmlquery"
)

func ParsePrivateLog(filePath string) (map[string]*FooocusMeta, error) {
	doc, err := htmlquery.LoadDoc(filePath)
	if err != nil {
		return nil, err
	}

	nodes, err := htmlquery.QueryAll(doc, "//div[@class='image-container']")
	if err != nil {
		return nil, err
	}

	var images map[string]*FooocusMeta = make(map[string]*FooocusMeta, len(nodes))

	for _, n := range nodes {
		img := htmlquery.FindOne(n, "//img")
		imgSrc := htmlquery.SelectAttr(img, "src")

		// Metadata is encoded in the onclick handler that allows
		// to copy the metadata to the clipboard.
		b := htmlquery.FindOne(n, "//button")
		bClick := htmlquery.SelectAttr(b, "onclick")

		stripLeft := "to_clipboard("
		stripRight := "')"
		clean := bClick[len(stripLeft)+1 : len(bClick)-len(stripRight)]
		cleanU, err := url.QueryUnescape(clean)
		if err != nil {
			return nil, err
		}

		// Parse metadata
		var metadata *FooocusMeta
		var legacyMetadata *FooocusMetaLegacy

		// Prefer modern format
		if err := json.Unmarshal([]byte(cleanU), &metadata); err == nil {
			slog.Debug("Found modern metadata format",
				"file", imgSrc, "data", cleanU)
			images[imgSrc] = metadata
		} else {
			// Fallback to legacy format
			if err := json.Unmarshal([]byte(cleanU), &legacyMetadata); err != nil {
				return images, fmt.Errorf("failed to read Fooocus parameters: %w", err)
			} else {
				slog.Debug("Found legacy metadata format",
					"file", imgSrc, "data", cleanU)
				images[imgSrc] = legacyMetadata.toFooocusMeta()
			}
		}
	}

	return images, nil
}
