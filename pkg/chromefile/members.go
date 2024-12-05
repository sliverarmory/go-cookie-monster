package chromefile

import (
	"path/filepath"
	"syscall"
)

type FileType int

const (
	// Windows API Constants
	PROCESS_DUP_HANDLE    = 0x0040
	DUPLICATE_SAME_ACCESS = 0x00000002
	FILE_TYPE_DISK        = 0x0001
)

const (
	// File Types
	Cookies FileType = iota
	LoginData
)

const (
	// Windows Error Codes
	ERROR_SHARING_VIOLATION syscall.Errno = 32
)

// FilePath contains information about a Chrome file
type FilePath struct {
	Type         FileType
	RelativePath string // Path relative to Chrome User Data directory
}

// FileData tracks data from either source
type FileData struct {
	Data       []byte
	FromHandle bool
	Size       int64
	Path       string
	Error      string // New field to store error messages
}

// ChromeFiles holds results of our file access attempts
type ChromeFiles struct {
	Cookies   FileData
	LoginData FileData
}

var chromePaths = map[FileType]FilePath{
	Cookies: {
		Type:         Cookies,
		RelativePath: filepath.Join("Default", "Network", "Cookies"),
	},
	LoginData: {
		Type:         LoginData,
		RelativePath: filepath.Join("Default", "Login Data"), // Fixed path
	},
}
