package processes

import "golang.org/x/sys/windows"

const (
	SystemHandleInformation = 0x10
	ObjectNameInformation   = 1
)

// Load the ntdll functions we need
var (
	ntdll                    = windows.NewLazySystemDLL("ntdll.dll")
	ntQuerySystemInformation = ntdll.NewProc("NtQuerySystemInformation")
	ntQueryObject            = ntdll.NewProc("NtQueryObject")
)

type SystemHandleEntry struct {
	ProcessId        uint32  // ULONG
	ObjectTypeNumber byte    // UCHAR
	Flags            byte    // UCHAR
	HandleValue      uint16  // USHORT
	Object           uintptr // PVOID
	GrantedAccess    uint32  // ACCESS_MASK
}

// Important: we need to tell Go to use the same memory layout as C
type SystemHandleInformationEx struct {
	NumberOfHandles uint32               // 4 bytes
	_               uint32               // 4 bytes padding
	Handles         [1]SystemHandleEntry // array of handles
}

type HandleCount struct {
	Count    uint32
	Reserved uint32
}

// Windows structure types
type UnicodeName struct { // Changed from UnicodeString to avoid conflicts
	Length        uint16
	MaximumLength uint16
	Buffer        *uint16
}

type ObjectName struct { // Changed from ObjectNameInformation to avoid conflicts
	Name UnicodeName
}
