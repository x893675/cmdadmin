package util


import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	errorsutil "k8s.io/apimachinery/pkg/util/errors"
)

const (
	// DefaultErrorExitCode defines exit the code for failed action generally
	DefaultErrorExitCode = 1
	// PreFlightExitCode defines exit the code for preflight checks
	PreFlightExitCode = 2
	// ValidationExitCode defines the exit code validation checks
	ValidationExitCode = 3
)

// fatal prints the message if set and then exits.
func fatal(msg string, code int) {
	if len(msg) > 0 {
		// add newline if needed
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}

		fmt.Fprint(os.Stderr, msg)
	}
	os.Exit(code)
}

// CheckErr prints a user friendly error to STDERR and exits with a non-zero
// exit code. Unrecognized errors will be printed with an "error: " prefix.
//
// This method is generic to the command in use and may be used by non-Kubectl
// commands.
func CheckErr(err error) {
	checkErr(err, fatal)
}

// preflightError allows us to know if the error is a preflight error or not
// defining the interface here avoids an import cycle of pulling in preflight into the util package
type preflightError interface {
	Preflight() bool
}

// checkErr formats a given error as a string and calls the passed handleErr
// func with that string and an exit code.
func checkErr(err error, handleErr func(string, int)) {

	var msg string
	if err != nil {
		msg = fmt.Sprintf("%s\nTo see the stack trace of this error execute with --v=5 or higher", err.Error())
		// check if the verbosity level in klog is high enough and print a stack trace.
		f := flag.CommandLine.Lookup("v")
		if f != nil {
			// assume that the "v" flag contains a parseable Int32 as per klog's "Level" type alias,
			// thus no error from ParseInt is handled here.
			if v, e := strconv.ParseInt(f.Value.String(), 10, 32); e == nil {
				// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-instrumentation/logging.md
				// klog.V(5) - Trace level verbosity
				if v > 4 {
					msg = fmt.Sprintf("%+v", err)
				}
			}
		}
	}

	switch err.(type) {
	case nil:
		return
	case preflightError:
		handleErr(msg, PreFlightExitCode)
	case errorsutil.Aggregate:
		handleErr(msg, ValidationExitCode)

	default:
		handleErr(msg, DefaultErrorExitCode)
	}
}

// FormatErrMsg returns a human-readable string describing the slice of errors passed to the function
func FormatErrMsg(errs []error) string {
	var errMsg string
	for _, err := range errs {
		errMsg = fmt.Sprintf("%s\t- %s\n", errMsg, err.Error())
	}
	return errMsg
}
