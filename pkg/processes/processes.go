package processes

import (
	"fmt"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Process represents a found Windows process
type Process struct {
	ID   uint32
	Name string
}

// FindProcess searches for processes by name and returns all matches
func FindProcess(name string) ([]Process, error) {
	// Get snapshot handle of all system processes
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to create process snapshot: %v", err)
	}
	defer windows.CloseHandle(snapshot)

	// Set up process entry structure
	var pe32 windows.ProcessEntry32
	pe32.Size = uint32(unsafe.Sizeof(pe32))

	// Get first process
	err = windows.Process32First(snapshot, &pe32)
	if err != nil {
		return nil, fmt.Errorf("failed to get first process: %v", err)
	}

	// Store matching processes
	var processes []Process

	// Iterate through processes
	for {
		if windows.UTF16ToString(pe32.ExeFile[:]) == name {
			processes = append(processes, Process{
				ID:   pe32.ProcessID,
				Name: name,
			})
		}

		err = windows.Process32Next(snapshot, &pe32)
		if err != nil {
			if err == syscall.ERROR_NO_MORE_FILES {
				break
			}
			return nil, fmt.Errorf("failed to get next process: %v", err)
		}
	}

	return processes, nil
}

// IsRunning is a convenience function that returns true if the process is found
func IsRunning(name string) (bool, error) {
	processes, err := FindProcess(name)
	if err != nil {
		return false, err
	}
	return len(processes) > 0, nil
}

func GetProcessHandles() ([]SystemHandleEntry, error) {
	//fmt.Println("[DEBUG] Starting getProcessHandles()")

	size := uint32(1024 * 1024 * 8)
	buffer := make([]byte, size)

	status, _, _ := ntQuerySystemInformation.Call(
		uintptr(SystemHandleInformation),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(size),
		uintptr(unsafe.Pointer(&size)))

	if status != 0 {
		return nil, fmt.Errorf("NtQuerySystemInformation failed with %x", status)
	}

	header := (*HandleCount)(unsafe.Pointer(&buffer[0]))

	handleOffset := unsafe.Sizeof(HandleCount{})

	handles := make([]SystemHandleEntry, 0)
	handleSize := unsafe.Sizeof(SystemHandleEntry{})

	for i := uint32(0); i < header.Count; i++ {
		entry := (*SystemHandleEntry)(unsafe.Pointer(&buffer[handleOffset+uintptr(i)*handleSize]))
		if entry.ProcessId > 0 && entry.ProcessId < 65535 { // Basic sanity check
			handles = append(handles, *entry)
		}
	}

	return handles, nil
}

func GetHandlePath(handle windows.Handle) (string, error) {
	// First call to get required buffer size
	var size uint32
	status, _, _ := ntQueryObject.Call(
		uintptr(handle),
		ObjectNameInformation,
		0,
		0,
		uintptr(unsafe.Pointer(&size)))

	if status != 0 && size == 0 {
		return "", fmt.Errorf("failed initial query with status: %x", status)
	}

	// Allocate buffer with some extra space
	size += 32 // Add a bit of extra space
	buffer := make([]byte, size)
	status, _, _ = ntQueryObject.Call(
		uintptr(handle),
		ObjectNameInformation,
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(size),
		uintptr(unsafe.Pointer(&size)))

	if status != 0 {
		return "", fmt.Errorf("NtQueryObject failed with status: %x", status)
	}

	// Get name information
	nameInfo := (*ObjectName)(unsafe.Pointer(&buffer[0]))
	if nameInfo.Name.Length == 0 {
		return "", fmt.Errorf("no name information available")
	}

	// Convert UTF-16 buffer to string
	bufferSize := nameInfo.Name.Length / 2
	utf16buf := make([]uint16, bufferSize)
	for i := range utf16buf {
		utf16buf[i] = *(*uint16)(unsafe.Pointer(
			uintptr(unsafe.Pointer(nameInfo.Name.Buffer)) + uintptr(i)*2))
	}

	return string(utf16.Decode(utf16buf)), nil
}
