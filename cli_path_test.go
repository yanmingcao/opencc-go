package opencc

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestCLIPathHandling tests path handling with Windows-style paths
func TestCLIPathHandling(t *testing.T) {
	tests := []struct {
		name       string
		configPath string
	}{
		{"forward slashes", "data/config/s2t.json"},
	}
	if runtime.GOOS == "windows" {
		tests = append(tests, struct {
			name       string
			configPath string
		}{"backslashes", "data\\config\\s2t.json"})
	}
	tests = append(tests, struct {
		name       string
		configPath string
	}{"absolute", filepath.Join(getWD(t), filepath.FromSlash("data/config/s2t.json"))})

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
	t.Helper()
	// Get test working directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	return wd
}
