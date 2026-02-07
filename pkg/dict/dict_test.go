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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrSingleValueDictEntry(t *testing.T) {
	entry := NewStrSingleValueDictEntry("简体", "繁體")

	assert.Equal(t, "简体", entry.Key())
	assert.Equal(t, []string{"繁體"}, entry.Values())
	assert.Equal(t, "繁體", entry.GetDefault())
	assert.Equal(t, 1, entry.NumValues())
	assert.Equal(t, "简体\t繁體", entry.ToString())
	assert.Equal(t, 6, entry.KeyLength())
}

func TestStrMultiValueDictEntry(t *testing.T) {
	entry := NewStrMultiValueDictEntry("发", []string{"髪", "發"})

	assert.Equal(t, "发", entry.Key())
	assert.Equal(t, []string{"髪", "發"}, entry.Values())
	assert.Equal(t, "髪", entry.GetDefault())
	assert.Equal(t, 2, entry.NumValues())
}

func TestNoValueDictEntry(t *testing.T) {
	entry := NewNoValueDictEntry("test")

	assert.Equal(t, "test", entry.Key())
	assert.Equal(t, []string{}, entry.Values())
	assert.Equal(t, "test", entry.GetDefault())
	assert.Equal(t, 0, entry.NumValues())
}

func TestLexicon(t *testing.T) {
	lexicon := NewLexicon()

	// Add in non-alphabetical order
	lexicon.Add(NewStrSingleValueDictEntry("c", "C"))
	lexicon.Add(NewStrSingleValueDictEntry("a", "A"))
	lexicon.Add(NewStrSingleValueDictEntry("b", "B"))

	assert.Equal(t, 3, lexicon.Len())

	// Before sorting, check it's not sorted
	assert.False(t, lexicon.IsSorted())

	// Sort it
	lexicon.Sort()
	assert.True(t, lexicon.IsSorted())

	// Check uniqueness
	assert.True(t, lexicon.IsUnique(nil))

	// Add duplicate
	lexicon.Add(NewStrSingleValueDictEntry("a", "X"))
	lexicon.Sort()
	var dupkey string
	assert.False(t, lexicon.IsUnique(&dupkey))
	assert.Equal(t, "a", dupkey)
}

func TestTextDict(t *testing.T) {
	lexicon := NewLexicon()
	lexicon.Add(NewStrSingleValueDictEntry("简化", "簡化"))
	lexicon.Add(NewStrSingleValueDictEntry("简体", "簡體"))
	lexicon.Add(NewStrSingleValueDictEntry("汉字", "漢字"))
	lexicon.Sort()

	d := NewTextDict(lexicon)

	// Test exact match
	entry := d.Match("简体")
	assert.NotNil(t, entry)
	assert.Equal(t, "簡體", entry.GetDefault())

	// Test non-existent entry
	entry = d.Match("不存在")
	assert.Nil(t, entry)

	// Test MatchPrefix
	entry = d.MatchPrefix("简体字")
	assert.NotNil(t, entry)
	assert.Equal(t, "简体", entry.Key())
	assert.Equal(t, "簡體", entry.GetDefault())

	// Test KeyMaxLength
	assert.Equal(t, 6, d.KeyMaxLength())
}

func TestTextDictMatchPrefix(t *testing.T) {
	lexicon := NewLexicon()
	lexicon.Add(NewStrSingleValueDictEntry("a", "A"))
	lexicon.Add(NewStrSingleValueDictEntry("ab", "AB"))
	lexicon.Add(NewStrSingleValueDictEntry("abc", "ABC"))
	lexicon.Sort()

	d := NewTextDict(lexicon)

	// Should find longest prefix "abc"
	entry := d.MatchPrefix("abcdef")
	assert.NotNil(t, entry)
	assert.Equal(t, "abc", entry.Key())
}

func TestDictGroup(t *testing.T) {
	// Create first dictionary
	lexicon1 := NewLexicon()
	lexicon1.Add(NewStrSingleValueDictEntry("a", "A1"))
	lexicon1.Sort()
	d1 := NewTextDict(lexicon1)

	// Create second dictionary (with different value for 'a')
	lexicon2 := NewLexicon()
	lexicon2.Add(NewStrSingleValueDictEntry("a", "A2"))
	lexicon2.Add(NewStrSingleValueDictEntry("b", "B2"))
	lexicon2.Sort()
	d2 := NewTextDict(lexicon2)

	// Create group
	group := NewDictGroup([]Dict{d1, d2})

	// Should find in first dict
	entry := group.Match("a")
	assert.NotNil(t, entry)
	assert.Equal(t, "A1", entry.GetDefault())

	// Should find in second dict (not in first)
	entry = group.Match("b")
	assert.NotNil(t, entry)
	assert.Equal(t, "B2", entry.GetDefault())

	// Test KeyMaxLength
	assert.Equal(t, 1, group.KeyMaxLength())
}
