package cookiemonster

import (
	"encoding/hex"
	"fmt"
	"go-cookie-monster/pkg/decrypt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func ProcessCookiesMode(key, databasePath string, databaseBytes []byte) {
	var (
		reader    *decrypt.DBReader
		err       error
		tempDbUse bool
	)

	fmt.Printf("\n[*] Attempting to decrypt cookies...\n")

	// ensure we have a key
	if key == "" {
		log.Fatalf("key is required for cookies mode")
	}

	// parse the key
	keyBytes, err := parseKey(key)
	if err != nil {
		log.Fatalf("error parsing key: %v", err)
	}

	// ensure we have a database path or bytes
	if databasePath == "" && len(databaseBytes) == 0 {
		log.Fatalf("database path is required for cookies mode")
	}

	// no database path but bytes, write the bytes to a tmep file
	if databasePath == "" && databaseBytes != nil {
		tmpFile, err := os.CreateTemp("", "")
		if err != nil {
			log.Fatalf("error creating temp file: %v", err)
		}
		//defer tmpFile.Close()

		if _, err := tmpFile.Write(databaseBytes); err != nil {
			log.Fatalf("error writing to temp file: %v", err)
		}
		tmpFile.Close()
		fmt.Println("[*] Wrote database bytes to temp file:", tmpFile.Name())

		databasePath = tmpFile.Name()
		tempDbUse = true
	}

	// read in the database file
	reader, err = decrypt.NewDBReader(databasePath)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	//defer reader.Close()

	// query the cookies
	rows, err := reader.QueryCookies()
	if err != nil {
		log.Fatalf("error querying cookies: %v", err)
	}
	//defer rows.Close()

	// extract the cookies
	extractor := &decrypt.CookieExtractor{Rows: rows}
	var cookies []decrypt.Cookie

	for rows.Next() {
		cookie, err := extractor.ExtractCookie(keyBytes)
		if err != nil {
			// print the key
			log.Printf("key: %s", key)
			log.Fatalf("error extracting cookie: %v", err)
		}
		cookies = append(cookies, *cookie)
	}

	formatter := &decrypt.JSONFormatter{Cookies: cookies}
	output, err := formatter.Format()
	if err != nil {
		log.Fatalf("error formatting cookies: %v", err)
	}

	fmt.Printf("[+] Decrypted %d cookies:\n", len(cookies))
	fmt.Println(output)

	// cleanup
	rows.Close()
	reader.Close()

	// if we created a temp file, remove it
	if tempDbUse {
		err = os.Remove(databasePath)
		if err != nil {
			log.Printf("[warning] error removing temp database file: %v", err)
		}
	}
}

func parseKey(keyStr string) ([]byte, error) {
	if len(keyStr) == 128 {
		keyStr = strings.ReplaceAll(keyStr, "\\x", "")
		decoded, err := hex.DecodeString(keyStr)
		if err != nil {
			return nil, fmt.Errorf("failed to decode hex string: %w", err)
		}
		return decoded, nil
	}
	return []byte(keyStr), nil
}
