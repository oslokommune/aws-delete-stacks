package cmd_test

import (
	"testing"
)

func TestCmdUsage(t *testing.T) {
	testCases := []*DeleteCmdTestCase{
		{
			Name:                   "Should print usage if no arguments are specified",
			Args:                   "",
			ExpectError:            true,
			ExpectedStderrContains: []string{"Error: validation failed: IncludeFilter: cannot be blank."},
		},
		{
			Name:             "Should print error if only --exclude flag is specified",
			Args:             "--exclude other-stack",
			ListStacksOutput: nil,
			ExpectError:      true,
			ExpectedStderrContains: []string{
				"Error: accepts 1 arg(s), received 0",
			},
		},
	}

	RunTestCases(t, testCases)
}
