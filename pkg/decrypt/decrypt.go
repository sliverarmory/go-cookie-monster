package decrypt

import (
	"database/sql"
)

func NewDBReader(path string) (*DBReader, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return &DBReader{db: db}, nil
}

func (r *DBReader) Close() error {
	return r.db.Close()
}

func (r *DBReader) QueryCookies() (*sql.Rows, error) {
	query := `SELECT name, encrypted_value, host_key, path, expires_utc, 
              is_secure, is_httponly, samesite 
              FROM cookies`
	return r.db.Query(query)
}

func (e *CookieExtractor) ExtractCookie(keyBytes []byte) (*Cookie, error) {
	var c Cookie
	var encryptedValue []byte
	var sameSiteInt int
	var expiresUtc int64

	err := e.Rows.Scan(&c.Name, &encryptedValue, &c.Domain, &c.Path,
		&expiresUtc, &c.Secure, &c.HTTPOnly, &sameSiteInt)
	if err != nil {
		return nil, err
	}

	d, err := decryptAESGCM(keyBytes, encryptedValue)
	if err != nil {
		return nil, err
	}

	c.Value = string(d)
	c.ExpirationDate = convertTimestamp(expiresUtc)
	c.SameSite = convertSameSite(sameSiteInt)
	setDefaultValues(&c)

	return &c, nil
}
