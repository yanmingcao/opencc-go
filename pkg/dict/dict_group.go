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
	"sort"
)

// DictGroup represents a group of dictionaries that are searched sequentially
type DictGroup struct {
	dicts []Dict
}

// NewDictGroup creates a new DictGroup from a slice of dictionaries
func NewDictGroup(dicts []Dict) *DictGroup {
	return &DictGroup{
		dicts: dicts,
	}
}

// MatchExact performs exact matching, searching each dictionary in order
func (g *DictGroup) Match(word string) DictEntry {
	for _, d := range g.dicts {
		if entry := d.Match(word); entry != nil {
			return entry
		}
	}
	return nil
}

// MatchPrefix finds the longest matching prefix across all dictionaries
func (g *DictGroup) MatchPrefix(word string) DictEntry {
	var bestEntry DictEntry
	bestLen := 0

	for _, d := range g.dicts {
		entry := d.MatchPrefix(word)
		if entry != nil && entry.KeyLength() > bestLen {
			bestEntry = entry
			bestLen = entry.KeyLength()
		}
	}

	return bestEntry
}

// MatchAllPrefixes finds all matching prefixes across all dictionaries
func (g *DictGroup) MatchAllPrefixes(word string) []DictEntry {
	allEntries := make([]DictEntry, 0)

	for _, d := range g.dicts {
		entries := d.MatchAllPrefixes(word)
		allEntries = append(allEntries, entries...)
	}

	// Sort by key length (descending) and then by key (ascending)
	sort.Slice(allEntries, func(i, j int) bool {
		if allEntries[i].KeyLength() != allEntries[j].KeyLength() {
			return allEntries[i].KeyLength() > allEntries[j].KeyLength()
		}
		return allEntries[i].Key() < allEntries[j].Key()
	})

	// Remove duplicates while preserving order
	uniqueEntries := make([]DictEntry, 0)
	seen := make(map[string]bool)
	for _, entry := range allEntries {
		key := entry.Key()
		if !seen[key] {
			seen[key] = true
			uniqueEntries = append(uniqueEntries, entry)
		}
	}

	return uniqueEntries
}

// KeyMaxLength returns the maximum key length across all dictionaries
func (g *DictGroup) KeyMaxLength() int {
	maxLen := 0
	for _, d := range g.dicts {
		if l := d.KeyMaxLength(); l > maxLen {
			maxLen = l
		}
	}
	return maxLen
}

// GetLexicon returns a merged lexicon from all dictionaries
func (g *DictGroup) GetLexicon() *Lexicon {
	lexicon := NewLexicon()

	for _, d := range g.dicts {
		dictLexicon := d.GetLexicon()
		for i := 0; i < dictLexicon.Len(); i++ {
			lexicon.Add(dictLexicon.At(i))
		}
	}

	lexicon.Sort()
	return lexicon
}

// Len returns the number of dictionaries in the group
func (g *DictGroup) Len() int {
	return len(g.dicts)
}

// At returns the dictionary at the given index
func (g *DictGroup) At(index int) Dict {
	if index < 0 || index >= len(g.dicts) {
		return nil
	}
	return g.dicts[index]
}

// Dicts returns the underlying dictionaries slice
func (g *DictGroup) Dicts() []Dict {
	return g.dicts
}
