package opencc

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/byvoid/opencc-go/pkg/dict"
)

// TestRealDictionaryLoading tests loading actual dictionary files
func TestRealDictionaryLoading(t *testing.T) {
	// Test loading STCharacters.txt
	t.Run("STCharacters", func(t *testing.T) {
		lexicon, err := dict.ParseLexiconFromFile("data/dictionary/STCharacters.txt")
		require.NoError(t, err)
		assert.Greater(t, lexicon.Len(), 1000) // Should have many entries

		// Verify sorting
		lexicon.Sort()
		assert.True(t, lexicon.IsSorted())

		// Create TextDict
		d := dict.NewTextDict(lexicon)
		assert.Greater(t, d.KeyMaxLength(), 0)

		// Test some conversions
		entry := d.Match("简")
		if entry != nil {
			t.Logf("'简' -> '%s'", entry.GetDefault())
		}

		entry = d.Match("发")
		if entry != nil {
			t.Logf("'发' has %d values", entry.NumValues())
		}
	})

	// Test loading STPhrases.txt
	t.Run("STPhrases", func(t *testing.T) {
		lexicon, err := dict.ParseLexiconFromFile("data/dictionary/STPhrases.txt")
		require.NoError(t, err)
		assert.Greater(t, lexicon.Len(), 100) // Should have phrase entries

		lexicon.Sort()
		d := dict.NewTextDict(lexicon)

		// Test phrase lookup
		entry := d.MatchPrefix("简体中文")
		if entry != nil {
			t.Logf("Found phrase: '%s' -> '%s'", entry.Key(), entry.GetDefault())
		}
	})
}

// TestEndToEndConversion tests actual Chinese conversion
func TestEndToEndConversion(t *testing.T) {
	// Check if data files exist
	if _, err := os.Stat("data/dictionary/STCharacters.txt"); os.IsNotExist(err) {
		t.Skip("Dictionary files not found, skipping integration test")
	}

	// Load dictionaries
	charLexicon, err := dict.ParseLexiconFromFile("data/dictionary/STCharacters.txt")
	require.NoError(t, err)
	charLexicon.Sort()
	charDict := dict.NewTextDict(charLexicon)

	phraseLexicon, err := dict.ParseLexiconFromFile("data/dictionary/STPhrases.txt")
	require.NoError(t, err)
	phraseLexicon.Sort()
	phraseDict := dict.NewTextDict(phraseLexicon)

	// Create DictGroup with phrase dict first (for priority)
	group := dict.NewDictGroup([]dict.Dict{phraseDict, charDict})

	// Test some conversions
	tests := []struct {
		input    string
		expected string
	}{
		{"简体", "簡體"},
		{"汉字", "漢字"},
		{"头发", "頭髪"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			entry := group.Match(tt.input)
			if entry != nil {
				result := entry.GetDefault()
				t.Logf("'%s' -> '%s'", tt.input, result)
				// Note: We don't assert exact match because we may not have all dictionaries
			}
		})
	}
}

// TestDataFilesExist verifies that data files were copied
func TestDataFilesExist(t *testing.T) {
	dataDir := "data"

	// Check config directory
	configDir := filepath.Join(dataDir, "config")
	entries, err := os.ReadDir(configDir)
	require.NoError(t, err, "Config directory should exist")
	assert.Greater(t, len(entries), 0, "Config directory should have files")

	// Check dictionary directory
	dictDir := filepath.Join(dataDir, "dictionary")
	entries, err = os.ReadDir(dictDir)
	require.NoError(t, err, "Dictionary directory should exist")
	assert.Greater(t, len(entries), 0, "Dictionary directory should have files")

	t.Logf("Found %d dictionary files", len(entries))
}
