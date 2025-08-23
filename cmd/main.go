package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	fooocusmeta "github.com/fkleon/fooocus-metadata"
	"github.com/fkleon/fooocus-metadata/fooocus"
	"github.com/fkleon/fooocus-metadata/fooocusplus"
	"github.com/fkleon/fooocus-metadata/ruinedfooocus"
)

func main() {

	var debug, verbose bool

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s <mode>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "mode: 'extract' or 'embed' (required)\n")
	}

	extractCmd := flag.NewFlagSet("extract", flag.ExitOnError)
	extractCmd.BoolVar(&verbose, "verbose", false, "enable verbose logging")
	extractCmd.BoolVar(&debug, "debug", false, "enable debug logging")
	extractCmd.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: extract [flags] <path>")
		extractCmd.PrintDefaults()
		fmt.Fprintf(os.Stderr, "path: The file to read metadata from (required)\n")
	}

	var embedType, embedIn, embedOut string
	embedCmd := flag.NewFlagSet("embed", flag.ExitOnError)
	embedCmd.BoolVar(&verbose, "verbose", false, "enable verbose logging")
	embedCmd.BoolVar(&debug, "debug", false, "enable debug logging")
	embedCmd.StringVar(&embedType, "type", "fooocus", "the type of metadata to embed (fooocus, fooocusplus, ruinedfooocus)")
	embedCmd.StringVar(&embedIn, "in", "", "the file to read imagedata from (optional)")
	embedCmd.StringVar(&embedOut, "out", "", "the file to write metadata to (required)")
	embedCmd.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: embed [flags] | echo '<meta>'")
		embedCmd.PrintDefaults()
		fmt.Fprintf(os.Stderr, "meta: The metadata to write in JSON format (stdin)\n")
	}

	if len(os.Args) < 2 {
		fmt.Println("expected 'extract' or 'embed' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "extract":
		extractCmd.Parse(os.Args[2:])
		setLogLevel(debug, verbose)

		path := extractCmd.Arg(0)

		if path == "" {
			extractCmd.Usage()
			os.Exit(1)
		}

		extract(path)

	case "embed":
		embedCmd.Parse(os.Args[2:])
		setLogLevel(debug, verbose)

		if embedOut == "" {
			embedCmd.Usage()
			os.Exit(1)
		}

		err := embed(embedType, embedIn, embedOut)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(2)
		}
		fmt.Printf("Metadata successfully embedded into %s\n", embedOut)
	default:
		flag.Usage()
	}
}

func extract(path string) {

	if metadata, err := fooocusmeta.ExtractFromFile(path); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(2)
	} else {
		out, err := json.MarshalIndent(metadata.Params.Raw(), "", "  ")
		if err == nil {
			fmt.Print(string(out))
		}
	}
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
		var metadata fooocus.Metadata
		err := json.NewDecoder(os.Stdin).Decode(&metadata)
		if err != nil {
			return fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		writer := fooocus.NewFooocusMetadataWriter()
		if in != "" {
			return writer.CopyWrite(source, target, metadata)
		} else {
			return writer.Write(target, metadata)
		}
	case "fooocusplus":
		var metadata fooocusplus.Metadata
		err := json.NewDecoder(os.Stdin).Decode(&metadata)
		if err != nil {
			return fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		writer := fooocusplus.NewFooocusPlusMetadataWriter()
		if in != "" {
			return writer.CopyWrite(source, target, metadata)
		} else {
			return writer.Write(target, metadata)
		}
	case "ruinedfooocus":
		var metadata ruinedfooocus.Metadata
		err := json.NewDecoder(os.Stdin).Decode(&metadata)
		if err != nil {
			return fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		writer := ruinedfooocus.NewRuinedFooocusMetadataWriter()
		if in != "" {
			return writer.CopyWrite(source, target, metadata)
		} else {
			return writer.Write(target, metadata)
		}
	default:
		fmt.Printf("Unknown type: %s\n", t)
		os.Exit(1)
	}

	return nil
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
