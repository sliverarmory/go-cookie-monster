package chromefile

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// getFileFromDisk attempts to read a file directly from the filesystem
func getFileFromDisk(filePath FilePath) ([]byte, string, error) {
	// Get user's Chrome path
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		return nil, "", fmt.Errorf("could not get LOCALAPPDATA environment variable")
	}
	//fmt.Printf("DEBUG: LOCALAPPDATA = %s\n", localAppData)

	// Build full path
	chromeUserData := filepath.Join(localAppData, "Google", "Chrome", "User Data")
	fullPath := filepath.Join(chromeUserData, filePath.RelativePath)
	//fmt.Printf("DEBUG: Attempting to read file: %s\n", fullPath)

	// Check if file exists first
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return nil, "", fmt.Errorf("file does not exist: %s", fullPath)
	}

	// Try to read file
	data, err := os.ReadFile(fullPath)
	if err != nil {
		// Check for specific error conditions
		if os.IsPermission(err) {
			return nil, "", fmt.Errorf("permission denied accessing %s", fullPath)
		}

		// On Windows, check for sharing violation or file in use
		if pathErr, ok := err.(*os.PathError); ok {
			// Windows ERROR_SHARING_VIOLATION
			if errno, ok := pathErr.Err.(syscall.Errno); ok {
				if errno == ERROR_SHARING_VIOLATION {
					return nil, "", fmt.Errorf("file in use by another process: %s", fullPath)
				}
				// Add specific error number for debugging
				return nil, "", fmt.Errorf("file access error (%d) for %s: %v", errno, fullPath, err)
			}
		}

		// Generic error case
		return nil, "", fmt.Errorf("failed to read file %s: %v", fullPath, err)
	}

	//fmt.Printf("DEBUG: Successfully read %d bytes from file\n", len(data))
	return data, fullPath, nil
}
