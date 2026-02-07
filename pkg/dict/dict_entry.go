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

// Package dict provides dictionary types and interfaces for OpenCC-Go
package dict

import (
	"strings"
)

// DictEntry interface represents a dictionary entry with a key and values
type DictEntry interface {
	// Key returns the dictionary entry key
	Key() string

	// Values returns all values for this entry
	Values() []string

	// GetDefault returns the default (first) value
	GetDefault() string

	// NumValues returns the number of values
	NumValues() int

	// ToString returns a string representation (for serialization)
	ToString() string

	// KeyLength returns the byte length of the key
	KeyLength() int

	// LessThan compares with another entry (by key)
	LessThan(other DictEntry) bool

	// Equals compares with another entry (by key)
	Equals(other DictEntry) bool
}

// NoValueDictEntry represents a dictionary entry with only a key (no value)
type NoValueDictEntry struct {
	key string
}

// NewNoValueDictEntry creates a new NoValueDictEntry
func NewNoValueDictEntry(key string) *NoValueDictEntry {
	return &NoValueDictEntry{key: key}
}

// Key returns the dictionary entry key
func (e *NoValueDictEntry) Key() string {
	return e.key
}

// Values returns empty slice (no values)
func (e *NoValueDictEntry) Values() []string {
	return []string{}
}

// GetDefault returns the key itself (no value, so return key)
func (e *NoValueDictEntry) GetDefault() string {
	return e.key
}

// NumValues returns 0 (no values)
func (e *NoValueDictEntry) NumValues() int {
	return 0
}

// ToString returns the key
func (e *NoValueDictEntry) ToString() string {
	return e.key
}

// KeyLength returns the byte length of the key
func (e *NoValueDictEntry) KeyLength() int {
	return len(e.key)
}

// LessThan compares with another entry (by key)
func (e *NoValueDictEntry) LessThan(other DictEntry) bool {
	return e.key < other.Key()
}

// Equals compares with another entry (by key)
func (e *NoValueDictEntry) Equals(other DictEntry) bool {
	return e.key == other.Key()
}

// SingleValueDictEntry interface for entries with exactly one value
type SingleValueDictEntry interface {
	DictEntry
	Value() string
}

// StrSingleValueDictEntry represents a dictionary entry with a single key-value pair
type StrSingleValueDictEntry struct {
	key   string
	value string
}

// NewStrSingleValueDictEntry creates a new StrSingleValueDictEntry
func NewStrSingleValueDictEntry(key, value string) *StrSingleValueDictEntry {
	return &StrSingleValueDictEntry{key: key, value: value}
}

// Key returns the dictionary entry key
func (e *StrSingleValueDictEntry) Key() string {
	return e.key
}

// Value returns the single value
func (e *StrSingleValueDictEntry) Value() string {
	return e.value
}

// Values returns a slice with the single value
func (e *StrSingleValueDictEntry) Values() []string {
	return []string{e.value}
}

// GetDefault returns the single value
func (e *StrSingleValueDictEntry) GetDefault() string {
	return e.value
}

// NumValues returns 1 (single value)
func (e *StrSingleValueDictEntry) NumValues() int {
	return 1
}

// ToString returns "key\tvalue"
func (e *StrSingleValueDictEntry) ToString() string {
	return e.key + "\t" + e.value
}

// KeyLength returns the byte length of the key
func (e *StrSingleValueDictEntry) KeyLength() int {
	return len(e.key)
}

// LessThan compares with another entry (by key)
func (e *StrSingleValueDictEntry) LessThan(other DictEntry) bool {
	return e.key < other.Key()
}

// Equals compares with another entry (by key)
func (e *StrSingleValueDictEntry) Equals(other DictEntry) bool {
	return e.key == other.Key()
}

// MultiValueDictEntry interface for entries with multiple values
type MultiValueDictEntry interface {
	DictEntry
}

// StrMultiValueDictEntry represents a dictionary entry with multiple values
type StrMultiValueDictEntry struct {
	key    string
	values []string
}

// NewStrMultiValueDictEntry creates a new StrMultiValueDictEntry
func NewStrMultiValueDictEntry(key string, values []string) *StrMultiValueDictEntry {
	if len(values) == 0 {
		return &StrMultiValueDictEntry{key: key, values: []string{}}
	}
	if len(values) == 1 {
		return &StrMultiValueDictEntry{key: key, values: values}
	}
	return &StrMultiValueDictEntry{key: key, values: values}
}

// Key returns the dictionary entry key
func (e *StrMultiValueDictEntry) Key() string {
	return e.key
}

// Values returns all values
func (e *StrMultiValueDictEntry) Values() []string {
	return e.values
}

// GetDefault returns the first value (or key if no values)
func (e *StrMultiValueDictEntry) GetDefault() string {
	if len(e.values) > 0 {
		return e.values[0]
	}
	return e.key
}

// NumValues returns the number of values
func (e *StrMultiValueDictEntry) NumValues() int {
	return len(e.values)
}

// ToString returns "key\tvalue1\tvalue2\t..."
func (e *StrMultiValueDictEntry) ToString() string {
	return e.key + "\t" + strings.Join(e.values, "\t")
}

// KeyLength returns the byte length of the key
func (e *StrMultiValueDictEntry) KeyLength() int {
	return len(e.key)
}

// LessThan compares with another entry (by key)
func (e *StrMultiValueDictEntry) LessThan(other DictEntry) bool {
	return e.key < other.Key()
}

// Equals compares with another entry (by key)
func (e *StrMultiValueDictEntry) Equals(other DictEntry) bool {
	return e.key == other.Key()
}

// DictEntryFactory provides factory methods for creating dictionary entries
type DictEntryFactory struct{}

// New creates a new DictEntry from a key only (no values)
func (f *DictEntryFactory) New(key string) DictEntry {
	return NewNoValueDictEntry(key)
}

// NewSingle creates a new DictEntry from a key and single value
func (f *DictEntryFactory) NewSingle(key, value string) DictEntry {
	return NewStrSingleValueDictEntry(key, value)
}

// NewMulti creates a new DictEntry from a key and multiple values
func (f *DictEntryFactory) NewMulti(key string, values []string) DictEntry {
	if len(values) == 0 {
		return NewNoValueDictEntry(key)
	}
	if len(values) == 1 {
		return NewStrSingleValueDictEntry(key, values[0])
	}
	return NewStrMultiValueDictEntry(key, values)
}

// NewFromEntry creates a copy of an existing entry
func (f *DictEntryFactory) NewFromEntry(entry DictEntry) DictEntry {
	switch e := entry.(type) {
	case *NoValueDictEntry:
		return NewNoValueDictEntry(e.Key())
	case *StrSingleValueDictEntry:
		return NewStrSingleValueDictEntry(e.Key(), e.Value())
	case *StrMultiValueDictEntry:
		return NewStrMultiValueDictEntry(e.Key(), e.Values())
	default:
		// For unknown types, create based on available methods
		if entry.NumValues() == 0 {
			return NewNoValueDictEntry(entry.Key())
		}
		if entry.NumValues() == 1 {
			return NewStrSingleValueDictEntry(entry.Key(), entry.Values()[0])
		}
		return NewStrMultiValueDictEntry(entry.Key(), entry.Values())
	}
}

// Common factory instance
var EntryFactory = &DictEntryFactory{}
