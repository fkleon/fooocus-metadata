// A command-line tool to embed image generation parameters
// into an image file.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/fkleon/fooocus-metadata/fooocus"
	"github.com/fkleon/fooocus-metadata/fooocusplus"
	"github.com/fkleon/fooocus-metadata/ruinedfooocus"
)

func main() {

	var debug, verbose bool
	var embedType, embedIn, embedOut string

	flag.BoolVar(&verbose, "verbose", false, "enable verbose logging")
	flag.BoolVar(&debug, "debug", false, "enable debug logging")
	flag.StringVar(&embedType, "type", "fooocus", "the type of metadata to embed (fooocus, fooocusplus, ruinedfooocus)")
	flag.StringVar(&embedIn, "in", "", "the file to read imagedata from (optional)")
	flag.StringVar(&embedOut, "out", "", "the file to write metadata to (required)")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: [flags] | echo '<meta>'")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "meta: The metadata to write in JSON format (stdin)\n")
	}

	flag.Parse()
	setLogLevel(debug, verbose)

	if embedOut == "" {
		flag.Usage()
		os.Exit(1)
	}

	err := embed(embedType, embedIn, embedOut)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(2)
	}
	fmt.Printf("Metadata successfully embedded into %s\n", embedOut)
}

func embed(t string, in string, out string) (err error) {

	var source, target *os.File

	if in != "" {
		source, err = os.Open(in)
		if err != nil {
			return fmt.Errorf("failed to open source file for writing: %w", err)
		}
		defer source.Close()
	}

	target, err = os.OpenFile(out, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open target file for writing: %w", err)
	}
	defer target.Close()

	switch t {
	case "fooocus":
		if metadata, err := readMetadataFromStdin[fooocus.Metadata](); err != nil {
			return fmt.Errorf("failed to unmarshal metadata: %w", err)
		} else {
			writer := fooocus.NewFooocusMetadataWriter()
			if in != "" {
				return writer.CopyWrite(source, target, metadata)
			} else {
				return writer.Write(target, metadata)
			}
		}
	case "fooocusplus":
		if metadata, err := readMetadataFromStdin[fooocusplus.Metadata](); err != nil {
			return fmt.Errorf("failed to unmarshal metadata: %w", err)
		} else {
			writer := fooocusplus.NewFooocusPlusMetadataWriter()
			if in != "" {
				return writer.CopyWrite(source, target, metadata)
			} else {
				return writer.Write(target, metadata)
			}
		}
	case "ruinedfooocus":
		if metadata, err := readMetadataFromStdin[ruinedfooocus.Metadata](); err != nil {
			return fmt.Errorf("failed to unmarshal metadata: %w", err)
		} else {
			writer := ruinedfooocus.NewRuinedFooocusMetadataWriter()
			if in != "" {
				return writer.CopyWrite(source, target, metadata)
			} else {
				return writer.Write(target, metadata)
			}
		}
	default:
		fmt.Printf("Unknown type: %s\n", t)
		os.Exit(1)
	}

	return nil
}

func readMetadataFromStdin[M fooocus.Metadata | fooocusplus.Metadata | ruinedfooocus.Metadata]() (metadata M, err error) {
	err = json.NewDecoder(os.Stdin).Decode(&metadata)
	return metadata, err
}

func setLogLevel(debug bool, verbose bool) {
	if debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else if verbose {
		slog.SetLogLoggerLevel(slog.LevelInfo)
	} else {
		slog.SetLogLoggerLevel(slog.LevelWarn)
	}
}
