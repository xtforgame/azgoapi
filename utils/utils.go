package utils

import (
	"fmt"
	"os"
	"runtime"
	// "time"
)

// FinalReport is the panic handler and provides the last report
func FinalReport() {
	if err := recover(); err != nil {
		trace := make([]byte, 1024)
		count := runtime.Stack(trace, true)
		errMsg := fmt.Sprintf("%s", err)
		fmt.Printf("Recover from panic: %s\n", errMsg)
		fmt.Printf("Stack of %d bytes: %s\n", count, string(trace[:count]))
		os.Exit(2)
	} else {
	}
}
