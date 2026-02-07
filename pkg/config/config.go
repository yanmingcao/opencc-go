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

package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Common errors
var (
	ErrInvalidConfig   = errors.New("invalid configuration")
	ErrMissingField    = errors.New("missing required field")
	ErrUnknownDictType = errors.New("unknown dictionary type")
	ErrUnknownSegType  = errors.New("unknown segmentation type")
)

// Config represents the top-level configuration
type Config struct {
	Name            string                  `json:"name"`
	Segmentation    *SegmentationConfig     `json:"segmentation"`
	ConversionChain []*ConversionStepConfig `json:"conversion_chain"`
}

// SegmentationConfig represents segmentation configuration
type SegmentationConfig struct {
	Type SegmentationType `json:"type"`
	Dict *DictConfig      `json:"dict"`
}

// SegmentationType represents the type of segmentation algorithm
type SegmentationType string

const (
	SegmentationTypeMMseg SegmentationType = "mmseg"
)

// DictConfig represents dictionary configuration
type DictConfig struct {
	Type  string        `json:"type"`
	File  string        `json:"file"`
	Dicts []*DictConfig `json:"dicts,omitempty"`
}

// ConversionStepConfig represents a single conversion step configuration
type ConversionStepConfig struct {
	Dict *DictConfig `json:"dict"`
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return LoadConfigFromData(data, filepath.Dir(filename))
}

// LoadConfigFromData loads configuration from JSON data
func LoadConfigFromData(data []byte, configDir string) (*Config, error) {
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Resolve relative paths
	if err := config.resolvePaths(configDir); err != nil {
		return nil, err
	}

	return &config, nil
}

// resolvePaths converts relative paths to absolute paths
func (c *Config) resolvePaths(configDir string) error {
	// Resolve segmentation dictionary path
	if c.Segmentation != nil && c.Segmentation.Dict != nil {
		if err := c.Segmentation.Dict.resolvePath(configDir); err != nil {
			return err
		}
	}

	// Resolve conversion chain dictionary paths
	for _, step := range c.ConversionChain {
		if step.Dict != nil {
			if err := step.Dict.resolvePath(configDir); err != nil {
				return err
			}
		}
	}

	return nil
}

// resolvePath resolves a single dictionary path
// Note: We no longer prepend configDir here because the dictionary loader
// uses search paths to find files. This allows dictionaries to be in a
// different directory (e.g., data/dictionary) than the config (e.g., data/config).
func (d *DictConfig) resolvePath(configDir string) error {
	// Keep filenames as-is (relative), let the loader search for them
	// Only resolve if it's already an absolute path (which shouldn't change)

	// Recursively resolve nested dictionaries
	for _, dict := range d.Dicts {
		if err := dict.resolvePath(configDir); err != nil {
			return err
		}
	}

	return nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Segmentation == nil {
		return ErrMissingField
	}

	if c.Segmentation.Dict == nil {
		return ErrMissingField
	}

	return nil
}
