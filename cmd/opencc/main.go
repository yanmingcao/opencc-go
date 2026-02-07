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
)

const (
	version          = "1.1.7-go"
	defaultConfigDir = "data/config"
)

// Built-in configuration mappings (short name -> config file)
var configMappings = map[string]string{
	"s2t":   "s2t.json",
	"t2s":   "t2s.json",
	"s2tw":  "s2tw.json",
	"tw2s":  "tw2s.json",
	"s2hk":  "s2hk.json",
	"hk2s":  "hk2s.json",
	"s2twp": "s2twp.json",
	"tw2sp": "tw2sp.json",
	"hk2t":  "hk2t.json",
	"t2hk":  "t2hk.json",
	"jp2t":  "jp2t.json",
	"t2jp":  "t2jp.json",
	"tw2t":  "tw2t.json",
	"t2tw":  "t2tw.json",
}

func main() {
	var (
		configFile  = flag.String("c", "", "Configuration file or short name (e.g., s2t, t2s)")
		configLong  = flag.String("config", "", "Configuration file or short name")
		inputFile   = flag.String("i", "", "Input file (default: stdin)")
		inputLong   = flag.String("input", "", "Input file (default: stdin)")
		outputFile  = flag.String("o", "", "Output file (default: stdout)")
		outputLong  = flag.String("output", "", "Output file (default: stdout)")
		showVersion = flag.Bool("v", false, "Show version")
		versionLong = flag.Bool("version", false, "Show version")
		showHelp    = flag.Bool("h", false, "Show help")
		helpLong    = flag.Bool("help", false, "Show help")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Open Chinese Convert (OpenCC) %s - Go Port\n\n", version)
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -c, --config <name|file>  Conversion preset or config file\n")
		fmt.Fprintf(os.Stderr, "  -i, --input <file>        Input file (default: stdin)\n")
		fmt.Fprintf(os.Stderr, "  -o, --output <file>       Output file (default: stdout)\n")
		fmt.Fprintf(os.Stderr, "  -v, --version             Show version information\n")
		fmt.Fprintf(os.Stderr, "  -h, --help                Show this help message\n")
		fmt.Fprintf(os.Stderr, "\nConversion Presets:\n")
		fmt.Fprintf(os.Stderr, "  s2t    Simplified → Traditional (Mainland)\n")
		fmt.Fprintf(os.Stderr, "  t2s    Traditional → Simplified (Mainland)\n")
		fmt.Fprintf(os.Stderr, "  s2tw   Simplified → Traditional (Taiwan)\n")
		fmt.Fprintf(os.Stderr, "  tw2s   Traditional → Simplified (Taiwan)\n")
		fmt.Fprintf(os.Stderr, "  s2hk   Simplified → Traditional (Hong Kong)\n")
		fmt.Fprintf(os.Stderr, "  hk2s   Traditional → Simplified (Hong Kong)\n")
		fmt.Fprintf(os.Stderr, "  s2twp  Simplified → Traditional (Taiwan, with phrases)\n")
		fmt.Fprintf(os.Stderr, "  tw2sp  Traditional → Simplified (Taiwan, with phrases)\n")
		fmt.Fprintf(os.Stderr, "  jp2t   Japanese Kanji → Traditional Chinese\n")
		fmt.Fprintf(os.Stderr, "  t2jp   Traditional Chinese → Japanese Kanji\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -c s2t\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -c s2t -i input.txt -o output.txt\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  echo \"汉字\" | %s -c s2t\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -c data/config/s2t.json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nFor custom configurations, provide the full path to a JSON file.\n")
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
		fmt.Printf("OpenCC %s (Go Port)\n", version)
		os.Exit(0)
	}

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *configFile == "" {
		fmt.Fprintf(os.Stderr, "Error: Configuration is required (-c or --config option)\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Resolve config name to file path
	configPath := resolveConfig(*configFile)
	if configPath == "" {
		fmt.Fprintf(os.Stderr, "Error: Cannot find configuration: %s\n", *configFile)
		os.Exit(1)
	}

	// Create converter
	converter, err := opencc.NewSimpleConverter(configPath)
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

// resolveConfig resolves a config name or file path to an actual config file path
func resolveConfig(name string) string {
	// Check if it's a short name (no path separator, no .json extension)
	if !strings.ContainsAny(name, "/\\") && !strings.HasSuffix(name, ".json") {
		if configFile, ok := configMappings[name]; ok {
			// Try to find the config in the default config directory
			configPath := filepath.Join(defaultConfigDir, configFile)
			if _, err := os.Stat(configPath); err == nil {
				return configPath
			}
		}
	}

	// Otherwise, treat it as a file path and search in usual locations
	return findConfigFile(name)
}

// findConfigFile searches for a configuration file in various locations
func findConfigFile(filename string) string {
	// If it's already a full path, check if it exists
	if filepath.IsAbs(filename) {
		if _, err := os.Stat(filename); err == nil {
			return filename
		}
		return ""
	}

	// Search paths
	searchPaths := []string{
		".",                          // Current directory
		defaultConfigDir,             // Default config directory
		"/usr/share/opencc",          // System-wide Linux
		"/usr/local/share/opencc",    // Local installation Linux
		"/opt/homebrew/share/opencc", // Homebrew macOS
		"C:\\Program Files\\OpenCC\\share\\opencc", // Windows
	}

	for _, path := range searchPaths {
		fullPath := filepath.Join(path, filename)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath
		}
	}

	return ""
}
