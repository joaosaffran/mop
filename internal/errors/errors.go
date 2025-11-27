package errors

import (
	"fmt"
	"os"
)

// ExitOnError prints an error message and exits if err is not nil
func ExitOnError(err error, message string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
		os.Exit(1)
	}
}

// ExitWithMessage prints a message and exits with code 1
func ExitWithMessage(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

// ExitOnErrorf prints a formatted error message and exits if err is not nil
func ExitOnErrorf(err error, format string, args ...interface{}) {
	if err != nil {
		msg := fmt.Sprintf(format, args...)
		fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
		os.Exit(1)
	}
}
