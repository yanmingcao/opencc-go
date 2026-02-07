/*
 * Open Chinese Convert
 *
 * Copyright 2010-2020 Carbo Kuo <byvoid@byvoid.com>
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

package dict

import (
	"bufio"
	"errors"
	"io"
	"os"
	"sort"
)

// Dict interface represents a dictionary with matching operations
type Dict interface {
	// MatchExact performs exact matching and returns the entry or nil
	Match(word string) DictEntry

	// MatchPrefix finds the longest matching prefix
	MatchPrefix(word string) DictEntry

	// MatchAllPrefixes finds all matching prefixes, sorted by length (descending)
	MatchAllPrefixes(word string) []DictEntry

	// KeyMaxLength returns the maximum key length in the dictionary
	KeyMaxLength() int

	// GetLexicon returns all entries in the dictionary
	GetLexicon() *Lexicon
}

// Optional type for representing nullable values
type Optional struct {
	value DictEntry
	valid bool
}

// OptionalNone returns an invalid Optional
func OptionalNone() Optional {
	return Optional{valid: false}
}

// OptionalSome returns a valid Optional with the given value
func OptionalSome(entry DictEntry) Optional {
	return Optional{value: entry, valid: true}
}

// IsPresent returns true if the Optional contains a value
func (o Optional) IsPresent() bool {
	return o.valid
}

// Get returns the value or panics if not present
func (o Optional) Get() DictEntry {
	if !o.valid {
		panic("called Get on None Optional")
	}
	return o.value
}

// OrElse returns the value or the given default
func (o Optional) OrElse(defaultValue DictEntry) DictEntry {
	if o.valid {
		return o.value
	}
	return defaultValue
}

// TextDict is a dictionary implementation using sorted text storage
type TextDict struct {
	maxLength int
	lexicon   *Lexicon
}

// NewTextDict creates a new TextDict from a lexicon
// The lexicon must be sorted
func NewTextDict(lexicon *Lexicon) *TextDict {
	maxLen := 0
	for i := 0; i < lexicon.Len(); i++ {
		entry := lexicon.At(i)
		if entry.KeyLength() > maxLen {
			maxLen = entry.KeyLength()
		}
	}

	return &TextDict{
		maxLength: maxLen,
		lexicon:   lexicon,
	}
}

// NewTextDictFromFile creates a TextDict from a text file
func NewTextDictFromFile(filename string) (*TextDict, error) {
	lexicon, err := ParseLexiconFromFile(filename)
	if err != nil {
		return nil, err
	}
	lexicon.Sort()
	return NewTextDict(lexicon), nil
}

// MatchExact performs exact matching
func (d *TextDict) Match(word string) DictEntry {
	// Binary search for exact match
	idx := sort.Search(d.lexicon.Len(), func(i int) bool {
		return d.lexicon.At(i).Key() >= word
	})

	if idx < d.lexicon.Len() && d.lexicon.At(idx).Key() == word {
		return d.lexicon.At(idx)
	}
	return nil
}

// MatchPrefix finds the longest matching prefix
func (d *TextDict) MatchPrefix(word string) DictEntry {
	maxLen := min(len(word), d.maxLength)

	// Search from longest to shortest
	for l := maxLen; l > 0; l-- {
		prefix := word[:l]
		idx := sort.Search(d.lexicon.Len(), func(i int) bool {
			return d.lexicon.At(i).Key() >= prefix
		})

		if idx < d.lexicon.Len() && d.lexicon.At(idx).Key() == prefix {
			return d.lexicon.At(idx)
		}
	}

	return nil
}

// MatchAllPrefixes finds all matching prefixes, sorted by length (descending)
func (d *TextDict) MatchAllPrefixes(word string) []DictEntry {
	maxLen := min(len(word), d.maxLength)
	var results []DictEntry

	// Collect all matching prefixes
	for l := 1; l <= maxLen; l++ {
		prefix := word[:l]
		idx := sort.Search(d.lexicon.Len(), func(i int) bool {
			return d.lexicon.At(i).Key() >= prefix
		})

		if idx < d.lexicon.Len() && d.lexicon.At(idx).Key() == prefix {
			results = append(results, d.lexicon.At(idx))
		}
	}

	return results
}

// KeyMaxLength returns the maximum key length
func (d *TextDict) KeyMaxLength() int {
	return d.maxLength
}

// GetLexicon returns the lexicon
func (d *TextDict) GetLexicon() *Lexicon {
	return d.lexicon
}

// SerializableDict interface for dictionary serialization
type SerializableDict interface {
	Dict
	SerializeToFile(filename string) error
	SerializeToWriter(writer io.Writer) error
}

// SerializeToFile serializes the dictionary to a file
func (d *TextDict) SerializeToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return d.SerializeToWriter(file)
}

// SerializeToWriter serializes the dictionary to a writer
func (d *TextDict) SerializeToWriter(writer io.Writer) error {
	bufWriter := bufio.NewWriter(writer)
	defer bufWriter.Flush()

	for i := 0; i < d.lexicon.Len(); i++ {
		entry := d.lexicon.At(i)
		_, err := bufWriter.WriteString(entry.ToString() + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// BinaryDict represents a binary serialized dictionary
type BinaryDict struct {
	Dict
}

// Common errors for dictionary operations
var (
	ErrInvalidFormat  = errors.New("invalid dictionary format")
	ErrInvalidHeader  = errors.New("invalid dictionary header")
	ErrInvalidVersion = errors.New("unsupported dictionary version")
)

// DartsDictHeader is the file header for Darts format
var DartsDictHeader = []byte("OPENCCDARTS1")

// MarisaDictHeader is the file header for Marisa format
var MarisaDictHeader = []byte("OPENCC_MARISA_0.2.5")
