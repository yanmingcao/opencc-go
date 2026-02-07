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

package segmentation

import (
	"github.com/yanmingcao/opencc-go/pkg/dict"
)

// Segmentation interface for text segmentation strategies
type Segmentation interface {
	// Segment performs segmentation on the input text
	Segment(text string) *Segments
}

// MaxMatchSegmentation implements forward maximum matching segmentation
type MaxMatchSegmentation struct {
	dict dict.Dict
}

// NewMaxMatchSegmentation creates a new MaxMatchSegmentation with the given dictionary
func NewMaxMatchSegmentation(d dict.Dict) *MaxMatchSegmentation {
	return &MaxMatchSegmentation{
		dict: d,
	}
}

// Segment performs forward maximum matching segmentation
func (s *MaxMatchSegmentation) Segment(text string) *Segments {
	segments := NewSegments()

	if len(text) == 0 {
		return segments
	}

	position := 0
	textLength := len(text)
	maxKeyLen := s.dict.KeyMaxLength()

	for position < textLength {
		// Try to find the longest match starting from current position
		match := s.findLongestMatch(text, position, maxKeyLen)

		if match != nil {
			// Found a match, add it and advance
			matchStr := match.Key()
			segments.AddUnmanaged(&matchStr)
			position += match.KeyLength()
		} else {
			// No match found, advance by one character
			// Find the length of the next UTF-8 character
			charLen := s.nextUTF8CharLength(text, position)
			if charLen == 0 {
				// Invalid UTF-8, advance by one byte
				charLen = 1
			}
			// Add the character as-is
			charStr := text[position : position+charLen]
			segments.AddManaged(charStr)
			position += charLen
		}
	}

	return segments
}

// findLongestMatch finds the longest matching prefix in the dictionary
func (s *MaxMatchSegmentation) findLongestMatch(text string, start int, maxLen int) dict.DictEntry {
	maxMatch := dict.DictEntry(nil)

	// Try lengths from max to 1
	for l := maxLen; l > 0; l-- {
		if start+l > len(text) {
			continue
		}

		prefix := text[start : start+l]
		match := s.dict.Match(prefix)

		if match != nil {
			maxMatch = match
			break // Found the longest, no need to continue
		}
	}

	return maxMatch
}

// nextUTF8CharLength returns the byte length of the next UTF-8 character
func (s *MaxMatchSegmentation) nextUTF8CharLength(text string, position int) int {
	if position >= len(text) {
		return 0
	}

	b := text[position]

	// Determine character length from first byte
	switch {
	case b < 0x80:
		return 1
	case b < 0xE0:
		return 2
	case b < 0xF0:
		return 3
	case b < 0xF8:
		return 4
	default:
		return 0 // Invalid
	}
}

// GetDict returns the dictionary used for segmentation
func (s *MaxMatchSegmentation) GetDict() dict.Dict {
	return s.dict
}

// SegmentationType represents the type of segmentation algorithm
type SegmentationType string

const (
	// SegmentationTypeMMseg is maximum forward matching
	SegmentationTypeMMseg SegmentationType = "mmseg"
)

// SegmentationConfig represents configuration for creating a segmentation
type SegmentationConfig struct {
	Type SegmentationType
	Dict dict.Dict
}

// NewSegmentationFromConfig creates a segmentation from configuration
func NewSegmentationFromConfig(config *SegmentationConfig) Segmentation {
	switch config.Type {
	case SegmentationTypeMMseg:
		return NewMaxMatchSegmentation(config.Dict)
	default:
		// Default to maximum matching
		return NewMaxMatchSegmentation(config.Dict)
	}
}

// CharactersSegmentation performs character-by-character segmentation
// This is useful when no dictionary-based segmentation is needed
type CharactersSegmentation struct{}

// NewCharactersSegmentation creates a new CharactersSegmentation
func NewCharactersSegmentation() *CharactersSegmentation {
	return &CharactersSegmentation{}
}

// Segment performs character-by-character segmentation
func (s *CharactersSegmentation) Segment(text string) *Segments {
	segments := NewSegments()

	position := 0
	for position < len(text) {
		charLen := s.nextUTF8CharLength(text, position)
		if charLen == 0 {
			// Invalid UTF-8, skip one byte
			charLen = 1
		}
		charStr := text[position : position+charLen]
		segments.AddManaged(charStr)
		position += charLen
	}

	return segments
}

// nextUTF8CharLength returns the byte length of the next UTF-8 character
func (s *CharactersSegmentation) nextUTF8CharLength(text string, position int) int {
	if position >= len(text) {
		return 0
	}

	b := text[position]

	switch {
	case b < 0x80:
		return 1
	case b < 0xE0:
		return 2
	case b < 0xF0:
		return 3
	case b < 0xF8:
		return 4
	default:
		return 0
	}
}
