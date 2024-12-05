package keys

import (
	"bytes"
	"errors"
)

// ExtractKey searches for the given pattern in the data and extracts the key.
// The key starts after the pattern and ends at the next double-quote.
func ExtractKey(data []byte, pattern string) (string, error) {
	// Find the start of the pattern in the data
	start := bytes.Index(data, []byte(pattern))
	if start == -1 {
		return "", errors.New("pattern not found")
	}

	// Move the start pointer past the pattern
	start += len(pattern)

	// Find the end of the key (next double-quote)
	end := bytes.IndexByte(data[start:], '"')
	if end == -1 {
		return "", errors.New("end of key not found")
	}

	// Extract the key from the data
	key := data[start : start+end]

	return string(key), nil
}
