package stdredir

import (
	"io"
	"os"
	"sync"
)

// STD out and error
var (
	originalStdout *os.File
	originalStderr *os.File
	stdoutWr       *os.File
	stderrWr       *os.File
	captureLock    sync.Mutex
)

// StartCapture begins redirecting stdout and stderr to the provided writer.
func StartCapture(out io.Writer) {
	captureLock.Lock()
	defer captureLock.Unlock()

	// Create os.File wrappers around out for stdout and stderr
	stdoutPipe, stdoutW, _ := os.Pipe()
	stderrPipe, stderrW, _ := os.Pipe()

	originalStdout = os.Stdout
	originalStderr = os.Stderr

	os.Stdout = stdoutW
	os.Stderr = stderrW

	// Copy output in separate goroutines to avoid blocking
	go func() {
		io.Copy(out, stdoutPipe)
	}()
	go func() {
		io.Copy(out, stderrPipe)
	}()

	stdoutWr = stdoutW
	stderrWr = stderrW
}

// StopCapture stops redirecting stdout and stderr and restores them.
func StopCapture() {
	captureLock.Lock()
	defer captureLock.Unlock()

	if stdoutWr != nil {
		stdoutWr.Close()
	}
	if stderrWr != nil {
		stderrWr.Close()
	}

	os.Stdout = originalStdout
	os.Stderr = originalStderr
}
