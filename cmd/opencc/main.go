/*
 * Open Chinese Convert
 *
 * Copyright 2010-2014 Carbo Kuo <byvoid@byvoid.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/byvoid/opencc-go"
	"github.com/byvoid/opencc-go/pkg/embeddata"
)

const (
	version = "1.1.7-go"
)

// Built-in configuration mappings (short name -> config file)
var configMappings = map[string]string{
	"s2t":   "s2t",
	"t2s":   "t2s",
	"s2tw":  "s2tw",
	"tw2s":  "tw2s",
	"s2hk":  "s2hk",
	"hk2s":  "hk2s",
	"s2twp": "s2twp",
	"tw2sp": "tw2sp",
	"hk2t":  "hk2t",
	"t2hk":  "t2hk",
	"jp2t":  "jp2t",
	"t2jp":  "t2jp",
	"tw2t":  "tw2t",
	"t2tw":  "t2tw",
}

func main() {
	var (
		configFile  = flag.String("c", "", "Conversion preset (e.g., s2t, t2s, s2tw)")
		configLong  = flag.String("config", "", "Conversion preset or config file")
		inputFile   = flag.String("i", "", "Input file (default: stdin)")
		inputLong   = flag.String("input", "", "Input file (default: stdin)")
		outputFile  = flag.String("o", "", "Output file (default: stdout)")
		outputLong  = flag.String("output", "", "Output file (default: stdout)")
		showVersion = flag.Bool("v", false, "Show version")
		versionLong = flag.Bool("version", false, "Show version")
		showHelp    = flag.Bool("h", false, "Show help")
		helpLong    = flag.Bool("help", false, "Show help")
		listConfigs = flag.Bool("list", false, "List all available conversion presets")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "OpenCC-Go %s - Chinese Conversion Tool\n\n", version)
		fmt.Fprintf(os.Stderr, "Usage: opencc -c <preset|config-file> [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -c, --config <preset|file>  Conversion preset (e.g., s2t) or config file path\n")
		fmt.Fprintf(os.Stderr, "  -i, --input <file>         Input file (default: stdin)\n")
		fmt.Fprintf(os.Stderr, "  -o, --output <file>       Output file (default: stdout)\n")
		fmt.Fprintf(os.Stderr, "  -v, --version              Show version\n")
		fmt.Fprintf(os.Stderr, "  -h, --help                 Show this help\n")
		fmt.Fprintf(os.Stderr, "  --list                     List all available presets\n")
		fmt.Fprintf(os.Stderr, "\nConversion Presets (embedded):\n")
		fmt.Fprintf(os.Stderr, "  s2t    Simplified → Traditional (Mainland China)\n")
		fmt.Fprintf(os.Stderr, "  t2s    Traditional → Simplified (Mainland China)\n")
		fmt.Fprintf(os.Stderr, "  s2tw   Simplified → Traditional (Taiwan)\n")
		fmt.Fprintf(os.Stderr, "  tw2s   Traditional → Simplified (Taiwan)\n")
		fmt.Fprintf(os.Stderr, "  s2hk   Simplified → Traditional (Hong Kong)\n")
		fmt.Fprintf(os.Stderr, "  hk2s   Traditional → Simplified (Hong Kong)\n")
		fmt.Fprintf(os.Stderr, "  s2twp  Simplified → Traditional (Taiwan, with phrases)\n")
		fmt.Fprintf(os.Stderr, "  tw2sp  Traditional → Simplified (Taiwan, with phrases)\n")
		fmt.Fprintf(os.Stderr, "  jp2t   Japanese Kanji → Traditional Chinese\n")
		fmt.Fprintf(os.Stderr, "  t2jp   Traditional Chinese → Japanese Kanji\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  opencc -c s2t -i input.txt -o output.txt\n")
		fmt.Fprintf(os.Stderr, "  echo \"汉字\" | opencc -c s2t\n")
		fmt.Fprintf(os.Stderr, "  echo \"汉字\" | opencc -c data/config/s2t.json\n")
		fmt.Fprintf(os.Stderr, "\nNote: Use presets (s2t, t2s, etc.) for quick conversions without external files.\n")
		fmt.Fprintf(os.Stderr, "      Or provide a config file path for custom configurations.\n")
	}

	flag.Parse()

	// Handle both short and long options
	if *versionLong {
		*showVersion = true
	}
	if *helpLong {
		*showHelp = true
	}
	if *configLong != "" {
		*configFile = *configLong
	}
	if *inputLong != "" {
		*inputFile = *inputLong
	}
	if *outputLong != "" {
		*outputFile = *outputLong
	}

	if *showVersion {
		fmt.Printf("OpenCC-Go %s\n", version)
		os.Exit(0)
	}

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *listConfigs {
		fmt.Println("Available conversion presets:")
		for _, config := range embeddata.ListConfigs() {
			fmt.Printf("  %s\n", config)
		}
		os.Exit(0)
	}

	if *configFile == "" {
		fmt.Fprintf(os.Stderr, "Error: Conversion preset is required (-c or --config)\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Resolve config name to config content
	configContent, err := resolveConfig(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Cannot find configuration: %s\n", *configFile)
		os.Exit(1)
	}

	// Create converter from embedded config
	converter, err := opencc.NewSimpleConverterFromData(configContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create converter: %v\n", err)
		os.Exit(1)
	}

	// Open input
	var input io.Reader
	if *inputFile == "" {
		input = os.Stdin
	} else {
		file, err := os.Open(*inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Cannot open input file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		input = file
	}

	// Open output
	var output io.Writer
	if *outputFile == "" {
		output = os.Stdout
	} else {
		file, err := os.Create(*outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Cannot create output file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		output = file
	}

	// Process input line by line
	scanner := bufio.NewScanner(input)
	writer := bufio.NewWriter(output)
	defer writer.Flush()

	for scanner.Scan() {
		line := scanner.Text()
		converted := converter.Convert(line)
		writer.WriteString(converted)
		writer.WriteByte('\n')
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
}

// resolveConfig resolves a config name to its JSON content
// It supports:
// 1. File paths (absolute or relative) - reads from disk
// 2. Embedded preset names (e.g., "s2t", "t2s") - uses embedded data
func resolveConfig(name string) ([]byte, error) {
	// First, check if it's a file path that exists
	if filepath.IsAbs(name) || strings.ContainsAny(name, "/\\") {
		if _, err := os.Stat(name); err == nil {
			// File exists, read it
			data, err := os.ReadFile(name)
			if err == nil {
				return data, nil
			}
		}
	}

	// Check if it's a short preset name (no path separator, no .json extension)
	if !strings.ContainsAny(name, "/\\") && !strings.HasSuffix(name, ".json") {
		if _, ok := configMappings[name]; ok {
			// Look up in built-in mappings
			if configContent, err := embeddata.GetConfig(name); err == nil {
				return configContent, nil
			}
		}
	}

	// Check if it exists as embedded config (with or without .json)
	if embeddata.ConfigExists(name) {
		return embeddata.GetConfig(name)
	}

	// Check for .json extension
	if !strings.HasSuffix(name, ".json") {
		if embeddata.ConfigExists(name + ".json") {
			return embeddata.GetConfig(name + ".json")
		}
	}

	return nil, os.ErrNotExist
}
