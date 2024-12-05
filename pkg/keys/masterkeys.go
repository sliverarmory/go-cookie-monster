package keys

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// DATA_BLOB structure for CryptUnprotectData
type DATA_BLOB struct {
	cbData uint32
	pbData *byte
}

// cryptUnprotectData is a wrapper for the Windows CryptUnprotectData function.
func cryptUnprotectData(data []byte) ([]byte, error) {
	var out DATA_BLOB
	in := DATA_BLOB{
		cbData: uint32(len(data)),
		pbData: &data[0],
	}

	cryptUnprotectDataProc := syscall.MustLoadDLL("crypt32.dll").MustFindProc("CryptUnprotectData")

	ret, _, err := cryptUnprotectDataProc.Call(
		uintptr(unsafe.Pointer(&in)),
		0,
		0,
		0,
		0,
		0,
		uintptr(unsafe.Pointer(&out)),
	)
	if ret == 0 {
		return nil, fmt.Errorf("CryptUnprotectData failed: %v", err)
	}

	// Copy the decrypted data into a Go slice
	decrypted := make([]byte, out.cbData)
	copy(decrypted, unsafe.Slice(out.pbData, out.cbData))

	// Free memory allocated by CryptUnprotectData
	windows.LocalFree(windows.Handle(unsafe.Pointer(out.pbData)))

	return decrypted, nil
}

// GetMasterKey decodes a base64-encoded key, decrypts it using CryptUnprotectData, and returns the decrypted key.
func GetMasterKey(key string) (string, error) {
	// Decode the base64-encoded key
	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode key: %v", err)
	}

	// Skip the first 5 bytes
	if len(decoded) <= 5 {
		return "", errors.New("decoded key is too short")
	}
	decoded = decoded[5:]

	// Decrypt the key using DPAPI
	decrypted, err := cryptUnprotectData(decoded)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt key: %v", err)
	}

	// Convert the decrypted key to a hex string
	var buffer bytes.Buffer
	for _, b := range decrypted {
		buffer.WriteString(fmt.Sprintf("\\x%02x", b))
	}

	return buffer.String(), nil
}
