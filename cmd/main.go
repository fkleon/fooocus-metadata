package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	fooocusmeta "github.com/fkleon/fooocus-metadata"
	"github.com/fkleon/fooocus-metadata/fooocus"
	_ "github.com/fkleon/fooocus-metadata/fooocusplus"
	_ "github.com/fkleon/fooocus-metadata/ruinedfooocus"
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

	embedCmd := flag.NewFlagSet("embed", flag.ExitOnError)
	embedCmd.BoolVar(&verbose, "verbose", false, "enable verbose logging")
	embedCmd.BoolVar(&debug, "debug", false, "enable debug logging")
	embedCmd.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: embed [flags] <in> <out> <meta>")
		embedCmd.PrintDefaults()
		fmt.Fprintf(os.Stderr, "in: The file to read imagedata from (optional)\n")
		fmt.Fprintf(os.Stderr, "out: The file to write metadata to (required)\n")
		fmt.Fprintf(os.Stderr, "meta: The metadata to write in JSON format\n")
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

		if metadata, err := fooocusmeta.ExtractFromFile(path); err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(2)
		} else {
			out, err := json.MarshalIndent(metadata.Params.Raw(), "", "  ")
			if err == nil {
				fmt.Print(string(out))
			}
		}

	case "embed":
		embedCmd.Parse(os.Args[2:])
		setLogLevel(debug, verbose)

		var in, out, meta string

		switch len(embedCmd.Args()) {
		case 2:
			out = embedCmd.Arg(0)
			meta = embedCmd.Arg(1)
		case 3:
			in = embedCmd.Arg(0)
			out = embedCmd.Arg(1)
			meta = embedCmd.Arg(2)
			fmt.Printf("niy %s", in)
		default:
			embedCmd.Usage()
			os.Exit(1)
		}

		if out == "" || meta == "" {
			embedCmd.Usage()
			os.Exit(1)
		}

		var metadata fooocus.Metadata
		json.Unmarshal([]byte(meta), &metadata)

		writer := fooocus.NewFooocusMetadataWriter()
		print(writer)
		// TODO
		/*
			if in != "" {
				writer.CopyWrite(in, out, metadata)
			} else {
				writer.Write(out, metadata)
			}

			if err := fooocus.EmbedIntoFile(out, metadata); err != nil {
				fmt.Printf("Error: %s\n", err)
				os.Exit(2)
			}
		*/
	default:
		flag.Usage()
	}
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
