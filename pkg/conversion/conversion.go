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

package conversion

import (
	"github.com/byvoid/opencc-go/pkg/dict"
	"github.com/byvoid/opencc-go/pkg/segmentation"
)

// Conversion performs a single conversion step using a dictionary
type Conversion struct {
	dict dict.Dict
}

// NewConversion creates a new Conversion with the given dictionary
func NewConversion(d dict.Dict) *Conversion {
	return &Conversion{
		dict: d,
	}
}

// Convert converts a single phrase
func (c *Conversion) Convert(phrase string) string {
	if len(phrase) == 0 {
		return phrase
	}

	entry := c.dict.Match(phrase)
	if entry != nil {
		return entry.GetDefault()
	}
	return phrase
}

// Convert converts segmented text
func (c *Conversion) ConvertSegments(segments *segmentation.Segments) *segmentation.Segments {
	result := segmentation.NewSegments()

	iterator := segments.Iterator()
	for iterator.Next() {
		segment := iterator.Value()
		converted := c.Convert(segment)
		result.AddManaged(converted)
	}

	return result
}

// GetDict returns the dictionary used for conversion
func (c *Conversion) GetDict() dict.Dict {
	return c.dict
}

// ConversionChain represents a chain of conversions applied in sequence
type ConversionChain struct {
	conversions []*Conversion
}

// NewConversionChain creates a new ConversionChain from a slice of conversions
func NewConversionChain(conversions []*Conversion) *ConversionChain {
	return &ConversionChain{
		conversions: conversions,
	}
}

// Convert applies all conversions in the chain
func (c *ConversionChain) Convert(segments *segmentation.Segments) *segmentation.Segments {
	result := segments

	for _, conversion := range c.conversions {
		result = conversion.ConvertSegments(result)
	}

	return result
}

// GetConversions returns the list of conversions
func (c *ConversionChain) GetConversions() []*Conversion {
	return c.conversions
}

// ConversionConfig represents configuration for a conversion step
type ConversionConfig struct {
	Dict dict.Dict
}

// NewConversionFromConfig creates a Conversion from configuration
func NewConversionFromConfig(config *ConversionConfig) *Conversion {
	return NewConversion(config.Dict)
}
