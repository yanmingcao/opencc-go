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

package dict

import (
	"bufio"
	"os"
	"sort"
	"strings"
)

// Lexicon represents a collection of dictionary entries
type Lexicon struct {
	entries []DictEntry
}

// NewLexicon creates a new empty lexicon
func NewLexicon() *Lexicon {
	return &Lexicon{
		entries: make([]DictEntry, 0),
	}
}

// NewLexiconFromEntries creates a lexicon from a slice of entries
func NewLexiconFromEntries(entries []DictEntry) *Lexicon {
	lexicon := &Lexicon{
		entries: make([]DictEntry, len(entries)),
	}
	copy(lexicon.entries, entries)
	return lexicon
}

// Add adds an entry to the lexicon
func (l *Lexicon) Add(entry DictEntry) {
	l.entries = append(l.entries, entry)
}

// AddEntry adds an entry pointer to the lexicon
func (l *Lexicon) AddEntry(entry *DictEntry) {
	l.entries = append(l.entries, *entry)
}

// Len returns the number of entries
func (l *Lexicon) Len() int {
	return len(l.entries)
}

// At returns the entry at the given index
func (l *Lexicon) At(index int) DictEntry {
	if index < 0 || index >= len(l.entries) {
		return nil
	}
	return l.entries[index]
}

// Iterator returns an iterator for the lexicon
func (l *Lexicon) Iterator() *LexiconIterator {
	return &LexiconIterator{
		lexicon: l,
		index:   0,
	}
}

// Sort sorts the lexicon by key
func (l *Lexicon) Sort() {
	sort.Slice(l.entries, func(i, j int) bool {
		return l.entries[i].Key() < l.entries[j].Key()
	})
}

// IsSorted checks if the lexicon is sorted by key
func (l *Lexicon) IsSorted() bool {
	for i := 1; i < len(l.entries); i++ {
		if l.entries[i].Key() < l.entries[i-1].Key() {
			return false
		}
	}
	return true
}

// IsUnique checks if every key is unique (after sorting)
// Returns true if unique, false otherwise
// If dupkey is provided, it will be set to the first duplicate key found
func (l *Lexicon) IsUnique(dupkey *string) bool {
	if len(l.entries) == 0 {
		return true
	}

	// Make a copy and sort it
	entriesCopy := make([]DictEntry, len(l.entries))
	copy(entriesCopy, l.entries)
	sort.Slice(entriesCopy, func(i, j int) bool {
		return entriesCopy[i].Key() < entriesCopy[j].Key()
	})

	// Check for duplicates
	for i := 1; i < len(entriesCopy); i++ {
		if entriesCopy[i].Key() == entriesCopy[i-1].Key() {
			if dupkey != nil {
				*dupkey = entriesCopy[i].Key()
			}
			return false
		}
	}

	return true
}

// ParseLexiconFromFile parses a lexicon from a text file
// Format: tab-separated key-value pairs, one per line
func ParseLexiconFromFile(filename string) (*Lexicon, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ParseLexiconFromReader(bufio.NewReader(file))
}

// ParseLexiconFromReader parses a lexicon from a reader
func ParseLexiconFromReader(reader *bufio.Reader) (*Lexicon, error) {
	lexicon := NewLexicon()

	lineNum := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err.Error() != "EOF" {
			// Check for EOF with data
			if len(line) == 0 {
				break
			}
			if err.Error() == "EOF" {
				break
			}
		}
		lineNum++

		// Remove line ending
		line = strings.TrimRight(line, "\r\n")

		// Skip empty lines and comments
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			if err != nil {
				break
			}
			continue
		}

		// Parse tab-separated values
		parts := strings.SplitN(line, "\t", 2)
		key := parts[0]

		var values []string
		if len(parts) > 1 {
			// Multiple values are space-separated (after the tab)
			values = strings.Fields(parts[1])
		}

		// Create entry
		entry := EntryFactory.NewMulti(key, values)
		lexicon.Add(entry)

		if err != nil {
			break
		}
	}

	return lexicon, nil
}

// LexiconIterator provides iteration over a lexicon
type LexiconIterator struct {
	lexicon *Lexicon
	index   int
}

// Next advances to the next entry
func (it *LexiconIterator) Next() bool {
	it.index++
	return it.index < it.lexicon.Len()
}

// Value returns the current entry
func (it *LexiconIterator) Value() DictEntry {
	return it.lexicon.At(it.index)
}

// Index returns the current index
func (it *LexiconIterator) Index() int {
	return it.index
}

// Entries returns the underlying entries slice
func (l *Lexicon) Entries() []DictEntry {
	return l.entries
}

// SetEntries sets the entries directly
func (l *Lexicon) SetEntries(entries []DictEntry) {
	l.entries = entries
}
