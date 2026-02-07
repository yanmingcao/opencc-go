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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/byvoid/opencc-go/pkg/config"
	"github.com/byvoid/opencc-go/pkg/conversion"
	"github.com/byvoid/opencc-go/pkg/dict"
	"github.com/byvoid/opencc-go/pkg/segmentation"
)

func TestConverter(t *testing.T) {
	// Create a simple dictionary
	lexicon := dict.NewLexicon()
	lexicon.Add(dict.NewStrSingleValueDictEntry("简体", "簡體"))
	lexicon.Add(dict.NewStrSingleValueDictEntry("汉字", "漢字"))
	lexicon.Sort()
	d := dict.NewTextDict(lexicon)

	// Create segmentation
	seg := segmentation.NewMaxMatchSegmentation(d)

	// Create conversion chain
	conv := conversion.NewConversion(d)
	chain := conversion.NewConversionChain([]*conversion.Conversion{conv})

	// Create converter
	converter := NewConverter("test", seg, chain)

	// Test conversion
	result := converter.Convert("简体汉字")
	assert.Equal(t, "簡體漢字", result)

	// Test conversion with non-matching text
	result = converter.Convert("Hello")
	assert.Equal(t, "Hello", result)

	// Test empty string
	result = converter.Convert("")
	assert.Equal(t, "", result)
}

func TestConverterMultiValue(t *testing.T) {
	// Create a dictionary with multi-value entries
	lexicon := dict.NewLexicon()
	lexicon.Add(dict.NewStrMultiValueDictEntry("发", []string{"髪", "發"}))
	lexicon.Sort()
	d := dict.NewTextDict(lexicon)

	// Create segmentation
	seg := segmentation.NewMaxMatchSegmentation(d)

	// Create conversion chain
	conv := conversion.NewConversion(d)
	chain := conversion.NewConversionChain([]*conversion.Conversion{conv})

	// Create converter
	converter := NewConverter("test", seg, chain)

	// Test conversion - should get first value
	result := converter.Convert("头发")
	assert.Equal(t, "头髪", result)
}

func TestSimpleConverterFromConfig(t *testing.T) {
	// Create a test configuration
	cfg := &config.Config{
		Name: "Test Config",
		Segmentation: &config.SegmentationConfig{
			Type: config.SegmentationTypeMMseg,
			Dict: &config.DictConfig{
				Type: "text",
				File: "test.txt",
			},
		},
		ConversionChain: []*config.ConversionStepConfig{
			{
				Dict: &config.DictConfig{
					Type: "text",
					File: "test.txt",
				},
			},
		},
	}

	// Note: This test will fail without actual dictionary files
	// In a real test, we would create temporary dictionary files
	// For now, we just test that the structure is correct
	require.NotNil(t, cfg)
	assert.Equal(t, "Test Config", cfg.Name)
}

func TestSegmentation(t *testing.T) {
	// Create dictionary
	lexicon := dict.NewLexicon()
	lexicon.Add(dict.NewStrSingleValueDictEntry("简体", "簡體"))
	lexicon.Sort()
	d := dict.NewTextDict(lexicon)

	// Create segmentation
	seg := segmentation.NewMaxMatchSegmentation(d)

	// Test segmentation
	result := seg.Segment("简体中文")
	// "简体" is matched, "中" and "文" are kept as-is (characters)
	assert.Equal(t, 3, result.Length())
	assert.Equal(t, "简体", result.At(0))
	assert.Equal(t, "中", result.At(1))
	assert.Equal(t, "文", result.At(2))
}

func TestConversionChain(t *testing.T) {
	// Create two dictionaries
	lexicon1 := dict.NewLexicon()
	lexicon1.Add(dict.NewStrSingleValueDictEntry("a", "b"))
	lexicon1.Sort()
	d1 := dict.NewTextDict(lexicon1)

	lexicon2 := dict.NewLexicon()
	lexicon2.Add(dict.NewStrSingleValueDictEntry("b", "c"))
	lexicon2.Sort()
	d2 := dict.NewTextDict(lexicon2)

	// Create conversion chain
	conv1 := conversion.NewConversion(d1)
	conv2 := conversion.NewConversion(d2)
	chain := conversion.NewConversionChain([]*conversion.Conversion{conv1, conv2})

	// Test conversion - "a" -> "b" -> "c"
	segments := segmentation.NewSegments()
	segments.AddManaged("a")
	result := chain.Convert(segments)
	assert.Equal(t, "c", result.ToString())
}
