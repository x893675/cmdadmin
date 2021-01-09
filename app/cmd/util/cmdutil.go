package util

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// SubCmdRunE returns a function that handles a case where a subcommand must be specified
// Without this callback, if a user runs just the command without a subcommand,
// or with an invalid subcommand, cobra will print usage information, but still exit cleanly.
// We want to return an error code in these cases so that the
// user knows that their command was invalid.
func SubCmdRunE(name string) func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.Errorf("missing subcommand; %q is not meant to be run on its own", name)
		}

		return errors.Errorf("invalid subcommand: %q", args[0])
	}
}

// ValidateExactArgNumber validates that the required top-level arguments are specified
func ValidateExactArgNumber(args []string, supportedArgs []string) error {
	lenSupported := len(supportedArgs)
	validArgs := 0
	// Disregard possible "" arguments; they are invalid
	for _, arg := range args {
		if len(arg) > 0 {
			validArgs++
		}
		// break early for too many arguments
		if validArgs > lenSupported {
			return errors.Errorf("too many arguments. Required arguments: %v", supportedArgs)
		}
	}

	if validArgs < lenSupported {
		return errors.Errorf("missing one or more required arguments. Required arguments: %v", supportedArgs)
	}
	return nil
}
