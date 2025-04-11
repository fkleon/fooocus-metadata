package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/fkleon/fooocus-metadata"
)

func main() {

	var logFile string

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [flags] <path>\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "path: The file to read metadata from (required)\n")
	}
	flag.StringVar(&logFile, "log-file", "", "private log file location, defaults to 'log.html' in the same folder as the file.")
	flag.Parse()

	filePath := flag.Arg(0)

	if filePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	image, err := fooocus.NewImageInfo(filePath)
	if err != nil {
		slog.Error("Failed to extract image information",
			"error", err)
		os.Exit(2)
	}

	if image.FooocusMetadata == nil {
		// Fallback to private log file.
		// Use default log file location if none is given.
		if logFile == "" {
			logFile = path.Join(path.Dir(filePath), "log.html")
		}
		slog.Info("Extracting metadata from private log file..",
			"filepath", logFile)
		ppl, metadataErr := fooocus.ParsePrivateLog(logFile)
		if metadataErr == nil {
			image.FooocusMetadata = ppl[image.Name()]
		} else {
			slog.Warn("Failed to extract metadata from private log file",
				"filepath", logFile,
				"error", metadataErr)
		}
	}

	if image.FooocusMetadata == nil {
		fmt.Println("Error: No Fooocus metadata found")
		os.Exit(2)
	}

	out, err := json.MarshalIndent(image.FooocusMetadata, "", "  ")
	if err == nil {
		fmt.Print(string(out))
	}
}
