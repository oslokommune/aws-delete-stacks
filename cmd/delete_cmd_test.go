package cmd_test

import (
	"github.com/oslokommune/aws-delete-stacks/cloudformation_api"
	"testing"
)

type DeleteCmdTestCase struct {
	name                   string
	args                   string
	stackSummaries         []*cloudformation_api.StackSummary
	expectError            bool
	expectedStdoutContains []string
	expectedStderrContains []string
}

func TestCmdUsage(t *testing.T) {
	testCases := []*DeleteCmdTestCase{
		{
			name:                   "Should print usage if no arguments are specified",
			args:                   "",
			expectError:            true,
			expectedStderrContains: []string{"Error: validation failed: IncludeFilter: cannot be blank."},
		},
		{
			name:           "Should print error if only --exclude flag is specified",
			args:           "--exclude other-stack",
			stackSummaries: nil,
			expectError:    true,
			expectedStderrContains: []string{
				"Error: accepts 1 arg(s), received 0",
			},
		},
	}

	runTestCases(t, testCases)
}
