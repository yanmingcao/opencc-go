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

package opencc

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/byvoid/opencc-go/pkg/config"
	"github.com/byvoid/opencc-go/pkg/conversion"
	"github.com/byvoid/opencc-go/pkg/dict"
	"github.com/byvoid/opencc-go/pkg/segmentation"
)

// Converter is the main controller for segmentation and conversion
type Converter struct {
	name            string
	segmentation    segmentation.Segmentation
	conversionChain *conversion.ConversionChain
}

// NewConverter creates a new Converter
func NewConverter(name string, seg segmentation.Segmentation, chain *conversion.ConversionChain) *Converter {
	return &Converter{
		name:            name,
		segmentation:    seg,
		conversionChain: chain,
	}
}

// Convert converts the input text
func (c *Converter) Convert(text string) string {
	if len(text) == 0 {
		return text
	}

	// Step 1: Segment the input text
	segments := c.segmentation.Segment(text)

	// Step 2: Apply conversion chain
	result := c.conversionChain.Convert(segments)

	// Step 3: Concatenate result
	return result.ToString()
}

// ConvertToBuffer converts text and writes to the provided buffer
// Returns the number of bytes written
func (c *Converter) ConvertToBuffer(input string, buffer []byte) int {
	result := c.Convert(input)
	if len(result) > len(buffer) {
		result = result[:len(buffer)]
	}
	copy(buffer, result)
	return len(result)
}

// GetSegmentation returns the segmentation object
func (c *Converter) GetSegmentation() segmentation.Segmentation {
	return c.segmentation
}

// GetConversionChain returns the conversion chain
func (c *Converter) GetConversionChain() *conversion.ConversionChain {
	return c.conversionChain
}

// GetName returns the converter name
func (c *Converter) GetName() string {
	return c.name
}

// SimpleConverter provides a simple high-level API
type SimpleConverter struct {
	converter *Converter
}

// NewSimpleConverter creates a SimpleConverter from a configuration file
func NewSimpleConverter(configFilename string, searchPaths ...string) (*SimpleConverter, error) {
	// Get the directory of the config file
	configDir := filepath.Dir(configFilename)
	if configDir == "" {
		configDir = "."
	}

	// Add config directory and dictionary paths to search paths
	allPaths := append([]string{configDir, filepath.Join(configDir, "..", "dictionary"), "data", "data/dictionary"}, searchPaths...)

	// Load configuration
	cfg, err := config.LoadConfig(configFilename)
	if err != nil {
		return nil, err
	}

	return NewSimpleConverterFromConfig(cfg, allPaths...)
}

// NewSimpleConverterFromConfig creates a SimpleConverter from a Config object
func NewSimpleConverterFromConfig(cfg *config.Config, searchPaths ...string) (*SimpleConverter, error) {
	// Build search paths
	paths := append([]string{"data", "data/dictionary"}, searchPaths...)

	// Create segmentation
	seg, err := createSegmentation(cfg.Segmentation, paths)
	if err != nil {
		return nil, err
	}

	// Create conversion chain
	chain, err := createConversionChain(cfg.ConversionChain, paths)
	if err != nil {
		return nil, err
	}

	converter := NewConverter(cfg.Name, seg, chain)
	return &SimpleConverter{converter: converter}, nil
}

// Convert converts the input text
func (s *SimpleConverter) Convert(text string) string {
	return s.converter.Convert(text)
}

// Convert converts a null-terminated C-style string
func (s *SimpleConverter) ConvertCString(input string) string {
	// Find null terminator
	for i, ch := range input {
		if ch == 0 { // null character
			return s.converter.Convert(input[:i])
		}
	}
	return s.converter.Convert(input)
}

// ConvertWithLength converts a string with the specified length
func (s *SimpleConverter) ConvertWithLength(input string, length int) string {
	if length <= 0 {
		return ""
	}
	if length >= len(input) {
		return s.converter.Convert(input)
	}
	return s.converter.Convert(input[:length])
}

// ConvertToBuffer converts text and writes to the provided buffer
// Returns the number of bytes written
func (s *SimpleConverter) ConvertToBuffer(input string, buffer []byte) int {
	return s.converter.ConvertToBuffer(input, buffer)
}

// ConvertToBufferWithLength converts text with the specified length and writes to buffer
func (s *SimpleConverter) ConvertToBufferWithLength(input string, length int, buffer []byte) int {
	text := s.ConvertWithLength(input, length)
	if len(text) > len(buffer) {
		text = text[:len(buffer)]
	}
	copy(buffer, text)
	return len(text)
}

// GetConverter returns the underlying Converter
func (s *SimpleConverter) GetConverter() *Converter {
	return s.converter
}

// createSegmentation creates a Segmentation from configuration
func createSegmentation(cfg *config.SegmentationConfig, searchPaths []string) (segmentation.Segmentation, error) {
	dict, err := loadDictFromConfig(cfg.Dict, searchPaths)
	if err != nil {
		return nil, err
	}

	switch cfg.Type {
	case config.SegmentationTypeMMseg:
		return segmentation.NewMaxMatchSegmentation(dict), nil
	default:
		return segmentation.NewMaxMatchSegmentation(dict), nil
	}
}

// createConversionChain creates a ConversionChain from configuration
func createConversionChain(steps []*config.ConversionStepConfig, searchPaths []string) (*conversion.ConversionChain, error) {
	conversions := make([]*conversion.Conversion, len(steps))

	for i, step := range steps {
		d, err := loadDictFromConfig(step.Dict, searchPaths)
		if err != nil {
			return nil, err
		}
		conversions[i] = conversion.NewConversion(d)
	}

	return conversion.NewConversionChain(conversions), nil
}

// loadDictFromConfig loads a dictionary from configuration
func loadDictFromConfig(cfg *config.DictConfig, searchPaths []string) (dict.Dict, error) {
	switch cfg.Type {
	case "group":
		// DictGroup - composite dictionary
		dicts := make([]dict.Dict, len(cfg.Dicts))
		for i, d := range cfg.Dicts {
			dict, err := loadDictFromConfig(d, searchPaths)
			if err != nil {
				return nil, err
			}
			dicts[i] = dict
		}
		return dict.NewDictGroup(dicts), nil

	case "text", "ocd":
		// TextDict or legacy format
		return loadTextDict(cfg.File, searchPaths)

	case "ocd2":
		// Default Marisa trie format
		return loadMarisaDict(cfg.File, searchPaths)

	default:
		return nil, config.ErrUnknownDictType
	}
}

// loadTextDict loads a text dictionary
func loadTextDict(filename string, searchPaths []string) (dict.Dict, error) {
	path := findFile(filename, searchPaths)
	if path == "" {
		return nil, fmt.Errorf("dictionary file not found: %s (searched in: %v)", filename, searchPaths)
	}

	lexicon, err := dict.ParseLexiconFromFile(path)
	if err != nil {
		return nil, err
	}

	lexicon.Sort()
	return dict.NewTextDict(lexicon), nil
}

// loadMarisaDict loads a Marisa trie dictionary (placeholder)
// In a full implementation, this would use the marisa-trie library
func loadMarisaDict(filename string, searchPaths []string) (dict.Dict, error) {
	path := findFile(filename, searchPaths)
	if path == "" {
		return nil, os.ErrNotExist
	}

	// For now, fall back to text dict
	// TODO: Implement proper Marisa trie support
	return loadTextDict(path, []string{filepath.Dir(path)})
}

// findFile searches for a file in the given paths
func findFile(filename string, searchPaths []string) string {
	if filepath.IsAbs(filename) {
		if _, err := os.Stat(filename); err == nil {
			return filename
		}
		return ""
	}

	for _, path := range searchPaths {
		fullPath := filepath.Join(path, filename)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath
		}
	}

	return ""
}
