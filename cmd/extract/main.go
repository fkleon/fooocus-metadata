// A command-line tool to extract image generation parameters
// from an image file or sidecar file.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/fkleon/fooocus-metadata/fooocus"
	_ "github.com/fkleon/fooocus-metadata/fooocusplus"
	_ "github.com/fkleon/fooocus-metadata/ruinedfooocus"
	_ "github.com/fkleon/fooocus-metadata/stablediffusion"
	_ "github.com/fkleon/fooocus-metadata/stablediffusioncpp"

	fooocusmeta "github.com/fkleon/fooocus-metadata"
)

func main() {
	var debug, verbose bool

	flag.BoolVar(&verbose, "verbose", false, "enable verbose logging")
	flag.BoolVar(&debug, "debug", false, "enable debug logging")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: [flags] <path>")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "path: The file to read metadata from (required)\n")
	}

	flag.Parse()
	setLogLevel(debug, verbose)

	path := flag.Arg(0)

	if path == "" {
		flag.Usage()
		os.Exit(1)
	}

	extract(path)
}

func extract(path string) {
	metadata, err := fooocusmeta.ExtractFromFile(path)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(2)
	}

	out, err := json.MarshalIndent(metadata.Params.Raw(), "", "  ")
	if err == nil {
		fmt.Print(string(out))
	} else {
		fmt.Printf("Error encoding metadata: %s\n", err)
		os.Exit(2)
	}
}

func setLogLevel(debug bool, verbose bool) {
	switch {
	case debug:
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case verbose:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	default:
		slog.SetLogLoggerLevel(slog.LevelWarn)
	}
}
