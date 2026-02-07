package opencc

import (
	"path/filepath"
	"testing"
)

// TestCLIPathHandling tests path handling with Windows-style paths
func TestCLIPathHandling(t *testing.T) {
	tests := []struct {
		name       string
		configPath string
	}{
		{"forward slashes", "data/config/s2t.json"},
		{"backslashes", "data\\config\\s2t.json"},
		{"absolute forward", filepath.Join(getWD(t), "data/config/s2t.json")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing config path: %s", tt.configPath)
			t.Logf("Config dir: %s", filepath.Dir(tt.configPath))

			converter, err := NewSimpleConverter(tt.configPath)
			if err != nil {
				t.Logf("Error: %v", err)
				t.Fail()
				return
			}

			result := converter.Convert("简体汉字")
			t.Logf("Result: %s", result)
		})
	}
}

func getWD(t *testing.T) string {
	// Get test working directory
	wd := ""
	// In test environment, use relative path
	return wd
}
