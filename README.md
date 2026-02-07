# OpenCC-Go

A **Go port** of OpenCC (Open Chinese Convert) - a conversion tool for Traditional/Simplified Chinese and regional variants.

This is a pure Go implementation of the original [OpenCC](https://github.com/BYVoid/OpenCC) project by BYVoid.

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

## Quick Start

### Build and Run

```bash
# Clone the repository
git clone https://github.com/byvoid/opencc-go.git
cd opencc-go

# Build the command-line tool
go build -o opencc ./cmd/opencc

# Convert text from stdin
echo "简体汉字" | ./opencc -c data/config/s2t.json
# Output: 簡體漢字

# Convert a file
./opencc -c data/config/s2t.json -i input.txt -o output.txt
```

### Or run without building

```bash
echo "简体汉字" | go run ./cmd/opencc -c data/config/s2t.json
```

## Installation

### As a Library

```bash
go get github.com/byvoid/opencc-go
```

### As a CLI Tool

```bash
go install github.com/byvoid/opencc-go/cmd/opencc@latest
```

Then use:
```bash
echo "简体汉字" | opencc -c /path/to/config/s2t.json
```

## Usage

### Library

```go
package main

import (
    "fmt"
    "github.com/byvoid/opencc-go"
)

func main() {
    // Create a converter from configuration file
    converter, err := opencc.NewSimpleConverter("data/config/s2t.json")
    if err != nil {
        panic(err)
    }
    
    // Convert text
    result := converter.Convert("简体汉字")
    fmt.Println(result)  // 簡體漢字
}
```

### Command-Line Tool

```bash
# Show help
./opencc -h

# Convert from stdin (default)
echo "简体汉字" | ./opencc -c data/config/s2t.json

# Convert from file to file
./opencc -c data/config/s2t.json -i input.txt -o output.txt

# Convert from file to stdout
./opencc -c data/config/s2t.json -i input.txt
```

### Available Configurations

| Config | Description |
|--------|-------------|
| `s2t.json` | Simplified Chinese to Traditional Chinese |
| `t2s.json` | Traditional Chinese to Simplified Chinese |
| `s2tw.json` | Simplified Chinese to Taiwan Traditional |
| `tw2s.json` | Taiwan Traditional to Simplified Chinese |
| `s2hk.json` | Simplified Chinese to Hong Kong Traditional |
| `hk2s.json` | Hong Kong Traditional to Simplified Chinese |
| `s2twp.json` | Simplified Chinese to Taiwan Traditional (with phrases) |
| `tw2sp.json` | Taiwan Traditional to Simplified Chinese (with phrases) |
| `jp2t.json` | Japanese Kanji to Traditional Chinese |
| `t2jp.json` | Traditional Chinese to Japanese Kanji |
| And more... | See `data/config/` directory |

## Introduction

OpenCC-Go is a pure Go implementation of the OpenCC project, providing conversion between Traditional Chinese, Simplified Chinese, and Japanese Kanji. It supports character-level and phrase-level conversion, character variant conversion, and regional idioms.

### Features

- **Pure Go**: No CGO dependencies, works on all Go-supported platforms
- **High Performance**: Efficient dictionary matching with trie data structures
- **Text Dictionary Format**: Uses plain text dictionaries (.txt) for portability
- **Flexible Configuration**: JSON-based configuration for custom conversion rules
- **Command-Line Tool**: Easy-to-use CLI for batch processing
- **Cross-Platform**: Windows, macOS, Linux compatible

## Architecture

OpenCC-Go follows a modular design:

```
┌─────────────────┐
│  SimpleConverter│  High-level API
├─────────────────┤
│    Converter    │  Main controller
├─────────────────┤
│  Segmentation   │  MaxMatchSegmentation
├─────────────────┤
│   Conversion    │  ConversionChain
├─────────────────┤
│  Dictionary     │  TextDict, DictGroup
└─────────────────┘
```

### Core Components

- **Dictionary System**: Interface with implementations for TextDict and DictGroup
- **Segmentation**: Maximum forward matching (mmseg) algorithm
- **Conversion**: Multi-stage conversion pipeline
- **Configuration**: JSON-based configuration loader

## Configuration Files

Configuration files are JSON-based and define:

1. **Segmentation**: How to split input text into segments (default: Maximum Forward Matching)
2. **Conversion Chain**: Ordered list of dictionary conversions

Example configuration (simplified to traditional):

```json
{
  "name": "Simplified to Traditional Chinese",
  "segmentation": {
    "type": "mmseg",
    "dict": {
      "type": "text",
      "file": "STPhrases.txt"
    }
  },
  "conversion_chain": [
    {
      "dict": {
        "type": "group",
        "dicts": [
          {"type": "text", "file": "STPhrases.txt"},
          {"type": "text", "file": "STCharacters.txt"}
        ]
      }
    }
  ]
}
```

## Dictionary Format

Dictionaries are tab-separated text files:

```
简体	簡體
汉字	漢字
```

Multiple values are supported for one-to-many mappings:

```
发	髪	發
```

## Testing

Run all tests:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

## Project Structure

```
opencc-go/
├── cmd/opencc/         # Command-line tool
├── pkg/
│   ├── utf8/           # UTF-8 utilities
│   ├── dict/           # Dictionary types and implementations
│   ├── segmentation/   # Text segmentation
│   ├── conversion/     # Conversion engine
│   └── config/         # Configuration loader
├── data/
│   ├── config/         # JSON configuration files
│   ├── dictionary/     # Text dictionary files
│   ├── icon/           # Project logo
│   └── scheme/         # Character disambiguation specs
└── demo/               # Demo program
```

## Differences from Original OpenCC

This Go port differs from the original C++ implementation in the following ways:

1. **Dictionary Format**: Uses text (.txt) dictionaries instead of binary (.ocd2) for simplicity and portability
2. **No Dictionary Compilation**: Reads text dictionaries directly without requiring compilation step
3. **Pure Go**: No CGO or external dependencies
4. **Simplified Architecture**: Focuses on core conversion functionality

## Roadmap

- [x] Core dictionary system
- [x] Maximum forward matching segmentation
- [x] Multi-stage conversion chain
- [x] JSON configuration support
- [x] Command-line tool
- [x] Comprehensive test coverage
- [ ] Performance benchmarks
- [ ] Optional: Binary dictionary format support

## License

Apache License 2.0

## Acknowledgments

This is a Go port of the original [OpenCC](https://github.com/BYVoid/OpenCC) project by BYVoid.
