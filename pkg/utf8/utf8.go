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

// Package utf8 provides UTF-8 string utilities for OpenCC-Go
package utf8

import (
	"errors"
	"unicode/utf8"
)

// InvalidUTF8 represents an invalid UTF-8 encoding error
var InvalidUTF8 = errors.New("invalid UTF-8")

// NextCharLength returns the byte length of the next UTF-8 character.
// Returns 0 if the byte sequence is invalid.
func NextCharLength(s string, index int) int {
	if index >= len(s) {
		return 0
	}
	// Check the first byte to determine character width
	c := s[index]
	switch {
	case c < 0x80:
		return 1
	case c < 0xE0:
		return 2
	case c < 0xF0:
		return 3
	case c < 0xF8:
		return 4
	case c < 0xFC:
		return 5
	case c < 0xFE:
		return 6
	default:
		return 0
	}
}

// NextCharLengthNoException returns the byte length of the next UTF-8 character.
// Returns 0 if the byte sequence is invalid, without throwing an exception.
func NextCharLengthNoException(s string, index int) int {
	length := NextCharLength(s, index)
	if length == 0 {
		return 0
	}
	// Validate the complete UTF-8 sequence
	if index+length > len(s) {
		return 0
	}
	return length
}

// PrevCharLength returns the byte length of the previous UTF-8 character
// before the given position.
func PrevCharLength(s string, index int) int {
	if index <= 0 {
		return 0
	}

	// Try different UTF-8 character lengths working backwards
	for length := 1; length <= 6; length++ {
		if index-length < 0 {
			break
		}
		// Check if a character of this length ends at index
		charStart := index - length
		if NextCharLength(s, charStart) == length {
			return length
		}
	}
	return 0
}

// NextChar returns a substring starting from the next UTF-8 character
func NextChar(s string, index int) string {
	if index >= len(s) {
		return ""
	}
	length := NextCharLength(s, index)
	if length == 0 {
		return s[index:]
	}
	return s[index+length:]
}

// PrevChar returns a substring ending at the previous UTF-8 character
func PrevChar(s string, index int) string {
	if index <= 0 {
		return ""
	}
	length := PrevCharLength(s, index)
	if length == 0 {
		return s[:index]
	}
	return s[:index-length]
}

// Length returns the number of UTF-8 characters in the string
func Length(s string) int {
	count := 0
	for i := 0; i < len(s); {
		length := NextCharLength(s, i)
		if length == 0 {
			// Invalid UTF-8, skip one byte
			i++
		} else {
			count++
			i += length
		}
	}
	return count
}

// FindNextInline finds the next occurrence of character ch in str,
// but only within the same line (stops at '\n', '\r', or end of string)
func FindNextInline(s string, ch byte) int {
	for i := 0; i < len(s); {
		length := NextCharLength(s, i)
		if length == 0 {
			break
		}
		// Check if this is a line ending
		if s[i] == '\n' || s[i] == '\r' {
			break
		}
		if s[i] == ch {
			return i
		}
		i += length
	}
	return -1
}

// IsLineEnding returns true if the character is a line ending
func IsLineEnding(ch byte) bool {
	return ch == '\n' || ch == '\r' || ch == 0
}

// FromSubstr creates a substring with the given byte length,
// ensuring it doesn't break a UTF-8 character
func FromSubstr(s string, start, length int) string {
	if start < 0 || start >= len(s) {
		return ""
	}

	// Find the actual end position without breaking UTF-8
	end := start
	maxEnd := min(start+length, len(s))
	for end < maxEnd {
		charLen := NextCharLength(s, end)
		if charLen == 0 {
			break
		}
		if end+charLen > maxEnd {
			break
		}
		end += charLen
	}

	return s[start:end]
}

// NotShorterThan returns true if the string is at least as long as
// the given byte length (without breaking UTF-8 characters)
func NotShorterThan(s string, byteLength int) bool {
	if byteLength <= 0 {
		return true
	}
	if byteLength > len(s) {
		return false
	}

	// Check if we can take byteLength bytes without breaking UTF-8
	end := 0
	for i := 0; i < byteLength && end < len(s); {
		charLen := NextCharLength(s, end)
		if charLen == 0 {
			return false
		}
		end += charLen
	}

	return end >= byteLength
}

// TruncateUTF8 truncates the string to at most maxByteLength bytes,
// without breaking any UTF-8 character
func TruncateUTF8(s string, maxByteLength int) string {
	if maxByteLength <= 0 {
		return ""
	}
	if maxByteLength >= len(s) {
		return s
	}

	end := 0
	for end < maxByteLength && end < len(s) {
		charLen := NextCharLength(s, end)
		if charLen == 0 {
			break
		}
		if end+charLen > maxByteLength {
			break
		}
		end += charLen
	}

	return s[:end]
}

// ReplaceAll replaces all occurrences of 'from' with 'to' in the string
func ReplaceAll(s string, from string, to string) string {
	if from == "" {
		return s
	}

	result := s
	for {
		idx := findSubstring(result, from)
		if idx == -1 {
			break
		}
		result = result[:idx] + to + result[idx+len(from):]
	}
	return result
}

// findSubstring is a helper function to find substring
func findSubstring(s string, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// Join joins a slice of strings with the given separator
func Join(strs []string, separator string) string {
	if len(strs) == 0 {
		return ""
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += separator + strs[i]
	}
	return result
}

// GetByteMap returns a map from UTF-8 character index to byte offset
func GetByteMap(s string) []int {
	utf8Len := Length(s)
	byteMap := make([]int, utf8Len)

	byteOffset := 0
	for i := 0; i < utf8Len; i++ {
		byteMap[i] = byteOffset
		charLen := NextCharLength(s, byteOffset)
		if charLen == 0 {
			break
		}
		byteOffset += charLen
	}

	return byteMap
}

// RuneToString converts a rune to string
func RuneToString(r rune) string {
	return string(r)
}

// ValidateUTF8 checks if the string contains valid UTF-8 encoding
func ValidateUTF8(s string) bool {
	return utf8.ValidString(s)
}

// ValidFirstByte checks if a byte can be the start of a valid UTF-8 character
func ValidFirstByte(b byte) bool {
	// Check if it's a valid first byte for UTF-8
	return b >= 0xC2 && b <= 0xF4
}
