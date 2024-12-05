package decrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
)

func decryptAESGCM(key []byte, encryptedData []byte) ([]byte, error) {
	nonce := encryptedData[3:15]
	ciphertext := encryptedData[15 : len(encryptedData)-16]
	tag := encryptedData[len(encryptedData)-16:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesGCM.Open(nil, nonce, append(ciphertext, tag...), nil)
	if err != nil {
		return nil, err
	}

	return plaintext[32:], nil
}

func convertTimestamp(expiresUtc int64) int64 {
	return expiresUtc/1000000 - 11644473600
}

func convertSameSite(sameSiteInt int) string {
	switch sameSiteInt {
	case 0:
		return "no_restriction"
	case 1:
		return "lax"
	case 2:
		return "strict"
	default:
		return "no_restriction"
	}
}

func setDefaultValues(c *Cookie) {
	c.HostOnly = false
	c.Session = false
	c.FirstPartyDomain = ""
	c.PartitionKey = nil
	c.StoreID = nil
}

func (f *JSONFormatter) Format() (string, error) {
	jsonData, err := json.MarshalIndent(f.Cookies, "", "    ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
