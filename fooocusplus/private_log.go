package fooocusplus

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
)

func ParsePrivateLog(filePath string) (map[string]Metadata, error) {
	doc, err := htmlquery.LoadDoc(filePath)
	if err != nil {
		return nil, err
	}

	// Check that Log file is compatible with this parser
	title, err := htmlquery.Query(doc, "//title")
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(htmlquery.InnerText(title), "Fooocus Log") {
		return nil, fmt.Errorf("file is not a Fooocus Plus private log: %s", filePath)
	}

	// Find all images in the log file

	nodes, err := htmlquery.QueryAll(doc, "//div[@class='image-container']")
	if err != nil {
		return nil, err
	}

	var images map[string]Metadata = make(map[string]Metadata, len(nodes))

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
		var metadata MetadataPrivateLog

		if err := json.Unmarshal([]byte(cleanU), &metadata); err == nil {
			if !strings.HasPrefix(metadata.Version, "FooocusPlus ") {
				continue
			}
			slog.Debug("Metadata in private log", "file", imgSrc)
			images[imgSrc] = metadata.toMetadata()
		}
	}

	return images, nil
}
