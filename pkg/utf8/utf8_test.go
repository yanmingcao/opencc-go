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

package utf8

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextCharLength(t *testing.T) {
	tests := []struct {
		input    string
		index    int
		expected int
	}{
		{"hello", 0, 1},
		{"汉字", 0, 3},
		{"汉字", 3, 3},
		{"", 0, 0},
		{"abc", 3, 0},
	}

	for _, tt := range tests {
		result := NextCharLength(tt.input, tt.index)
		assert.Equal(t, tt.expected, result, "NextCharLength(%s, %d)", tt.input, tt.index)
	}
}

func TestLength(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"hello", 5},
		{"汉字", 2},
		{"", 0},
		{"Hello 世界", 8},
	}

	for _, tt := range tests {
		result := Length(tt.input)
		assert.Equal(t, tt.expected, result, "Length(%s)", tt.input)
	}
}

func TestTruncateUTF8(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"hello", 3, "hel"},
		{"汉字", 3, "汉"},
		{"汉字", 6, "汉字"},
		{"", 10, ""},
	}

	for _, tt := range tests {
		result := TruncateUTF8(tt.input, tt.maxLen)
		assert.Equal(t, tt.expected, result, "TruncateUTF8(%s, %d)", tt.input, tt.maxLen)
	}
}

func TestReplaceAll(t *testing.T) {
	tests := []struct {
		input    string
		from     string
		to       string
		expected string
	}{
		{"hello world", "world", "universe", "hello universe"},
		{"aaa", "a", "b", "bbb"},
		{"hello", "x", "y", "hello"},
		{"", "a", "b", ""},
	}

	for _, tt := range tests {
		result := ReplaceAll(tt.input, tt.from, tt.to)
		assert.Equal(t, tt.expected, result, "ReplaceAll(%s, %s, %s)", tt.input, tt.from, tt.to)
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		strs      []string
		separator string
		expected  string
	}{
		{[]string{"a", "b", "c"}, ",", "a,b,c"},
		{[]string{"a"}, ",", "a"},
		{[]string{}, ",", ""},
	}

	for _, tt := range tests {
		result := Join(tt.strs, tt.separator)
		assert.Equal(t, tt.expected, result, "Join(%v, %s)", tt.strs, tt.separator)
	}
}

func TestValidateUTF8(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"hello", true},
		{"汉字", true},
		{"", true},
		{string([]byte{0xff, 0xfe}), false},
	}

	for _, tt := range tests {
		result := ValidateUTF8(tt.input)
		assert.Equal(t, tt.expected, result, "ValidateUTF8(%s)", tt.input)
	}
}
