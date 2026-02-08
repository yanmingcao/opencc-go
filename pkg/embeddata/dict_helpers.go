package embeddata

import (
	"errors"
	"strings"
)

// ErrDictNotFound is returned when a dictionary is not found.
var ErrDictNotFound = errors.New("dictionary not found")

// GetDict returns the dictionary content for the given name.
// The name can be with or without the .txt extension.
func GetDict(name string) ([]byte, error) {
	baseName := strings.TrimSuffix(name, ".txt")
	if content, ok := EmbeddedDict[baseName]; ok {
		return []byte(content), nil
	}
	return nil, ErrDictNotFound
}

// DictExists returns true if a dictionary with the given name exists.
// The name can be with or without the .txt extension.
func DictExists(name string) bool {
	_, err := GetDict(name)
	return err == nil
}
