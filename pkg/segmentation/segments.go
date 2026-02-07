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

// Package segmentation provides text segmentation for OpenCC-Go
package segmentation

import (
	"strings"
)

// Segments represents segmented text
type Segments struct {
	// managed contains string copies
	managed []string
	// unmanaged contains pointers to external strings
	unmanaged []*string
	// indexes maps position to (source, isManaged)
	indexes []struct {
		source  int
		managed bool
	}
}

// NewSegments creates a new empty Segments
func NewSegments() *Segments {
	return &Segments{
		managed:   make([]string, 0),
		unmanaged: make([]*string, 0),
		indexes: make([]struct {
			source  int
			managed bool
		}, 0),
	}
}

// NewSegmentsFromStrings creates Segments from a slice of strings
func NewSegmentsFromStrings(strs []string) *Segments {
	segments := NewSegments()
	for _, s := range strs {
		segments.AddManaged(s)
	}
	return segments
}

// AddManaged adds a managed (owned) string segment
func (s *Segments) AddManaged(str string) {
	s.indexes = append(s.indexes, struct {
		source  int
		managed bool
	}{source: len(s.managed), managed: true})
	s.managed = append(s.managed, str)
}

// AddUnmanaged adds an unmanaged (borrowed) string segment
func (s *Segments) AddUnmanaged(str *string) {
	s.indexes = append(s.indexes, struct {
		source  int
		managed bool
	}{source: len(s.unmanaged), managed: false})
	s.unmanaged = append(s.unmanaged, str)
}

// AddString adds a string as either managed or unmanaged
func (s *Segments) AddString(str string, managed bool) {
	if managed {
		s.AddManaged(str)
	} else {
		// For unmanaged, we need to store a pointer
		s.AddUnmanaged(&str)
	}
}

// At returns the segment at the given position
func (s *Segments) At(pos int) string {
	if pos < 0 || pos >= len(s.indexes) {
		return ""
	}
	idx := s.indexes[pos]
	if idx.managed {
		return s.managed[idx.source]
	}
	return *(s.unmanaged[idx.source])
}

// Length returns the number of segments
func (s *Segments) Length() int {
	return len(s.indexes)
}

// ToString concatenates all segments into a single string
func (s *Segments) ToString() string {
	var builder strings.Builder
	for i := 0; i < s.Length(); i++ {
		builder.WriteString(s.At(i))
	}
	return builder.String()
}

// Iterator provides iteration over segments
type SegmentsIterator struct {
	segments *Segments
	position int
}

// Iterator returns an iterator for the segments
func (s *Segments) Iterator() *SegmentsIterator {
	return &SegmentsIterator{
		segments: s,
		position: -1,
	}
}

// Next advances to the next segment
func (it *SegmentsIterator) Next() bool {
	it.position++
	return it.position < it.segments.Length()
}

// Value returns the current segment
func (it *SegmentsIterator) Value() string {
	return it.segments.At(it.position)
}

// Position returns the current position
func (it *SegmentsIterator) Position() int {
	return it.position
}

// Managed returns a slice of all managed strings
func (s *Segments) Managed() []string {
	return s.managed
}

// Unmanaged returns a slice of all unmanaged string pointers
func (s *Segments) Unmanaged() []*string {
	return s.unmanaged
}

// Indexes returns the index mapping
func (s *Segments) Indexes() []struct {
	source  int
	managed bool
} {
	return s.indexes
}
