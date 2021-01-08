package util

import (
	"testing"
)

func TestValidateExactArgNumber(t *testing.T) {
	var tests = []struct {
		name                string
		args, supportedArgs []string
		expectedErr         bool
	}{
		{
			name:          "one arg given and one arg expected",
			args:          []string{"my-node-1234"},
			supportedArgs: []string{"node-name"},
			expectedErr:   false,
		},
		{
			name:          "two args given and two args expected",
			args:          []string{"my-node-1234", "foo"},
			supportedArgs: []string{"node-name", "second-toplevel-arg"},
			expectedErr:   false,
		},
		{
			name:          "too few supplied args",
			args:          []string{},
			supportedArgs: []string{"node-name"},
			expectedErr:   true,
		},
		{
			name:          "too few non-empty args",
			args:          []string{""},
			supportedArgs: []string{"node-name"},
			expectedErr:   true,
		},
		{
			name:          "too many args",
			args:          []string{"my-node-1234", "foo"},
			supportedArgs: []string{"node-name"},
			expectedErr:   true,
		},
	}
	for _, rt := range tests {
		t.Run(rt.name, func(t *testing.T) {
			actual := ValidateExactArgNumber(rt.args, rt.supportedArgs)
			if (actual != nil) != rt.expectedErr {
				t.Errorf(
					"failed ValidateExactArgNumber:\n\texpected error: %t\n\t  actual error: %t",
					rt.expectedErr,
					(actual != nil),
				)
			}
		})
	}
}