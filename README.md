# OpenCC-Go

A Go implementation of OpenCC (Open Chinese Convert) - a conversion tool for Traditional/Simplified Chinese and regional variants.

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

## Introduction

OpenCC-Go is a pure Go implementation of the OpenCC project, providing conversion between Traditional Chinese, Simplified Chinese, and Japanese Kanji. It supports character-level and phrase-level conversion, character variant conversion, and regional idioms.

### Features

- **Pure Go**: No CGO dependencies, works on all Go-supported platforms
- **High Performance**: Efficient dictionary matching with trie data structures
- **Multiple Formats**: Support for text (.txt), OCD (legacy), and OCD2 (default) dictionary formats
- **Flexible Configuration**: JSON-based configuration for custom conversion rules
- **Command-Line Tool**: Easy-to-use CLI for batch processing

## Installation

```bash
go get github.com/byvoid/opencc-go
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
    converter, err := opencc.NewSimpleConverter("s2t.json")
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
# Build the command-line tool
go build ./cmd/opencc

# Convert text from stdin
echo "简体汉字" | ./opencc -c s2t.json

# Convert a file
./opencc -c s2t.json -i input.txt -o output.txt
```

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

- **Dictionary System**: Interface with implementations for TextDict, MarisaDict, DartsDict, and DictGroup
- **Segmentation**: Maximum forward matching (mmseg) algorithm
- **Conversion**: Multi-stage conversion pipeline
- **Configuration**: JSON-based configuration loader

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
├── data/               # Dictionary and config files
└── test/               # Test cases
```

## Roadmap

- [x] Core dictionary system
- [x] Maximum forward matching segmentation
- [x] Multi-stage conversion chain
- [x] JSON configuration support
- [x] Command-line tool
- [ ] Marisa trie integration for OCD2 format
- [ ] Darts double-array trie for OCD format
- [ ] Dictionary compilation tool
- [ ] Performance benchmarks
- [ ] More comprehensive test coverage

## License

Apache License 2.0

## Acknowledgments

This is a Go port of the original [OpenCC](https://github.com/BYVoid/OpenCC) project by BYVoid.
