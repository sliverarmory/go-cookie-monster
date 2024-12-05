package decrypt

import "database/sql"

type Cookie struct {
	Name             string      `json:"name"`
	Value            string      `json:"value"`
	Domain           string      `json:"domain"`
	HostOnly         bool        `json:"hostOnly"`
	Path             string      `json:"path"`
	Secure           bool        `json:"secure"`
	HTTPOnly         bool        `json:"httpOnly"`
	SameSite         string      `json:"sameSite"`
	Session          bool        `json:"session"`
	FirstPartyDomain string      `json:"firstPartyDomain"`
	PartitionKey     interface{} `json:"partitionKey"`
	ExpirationDate   int64       `json:"expirationDate"`
	StoreID          interface{} `json:"storeId"`
}

type DBReader struct {
	db *sql.DB
}

type CookieExtractor struct {
	Rows *sql.Rows
}

type JSONFormatter struct {
	Cookies []Cookie
}
