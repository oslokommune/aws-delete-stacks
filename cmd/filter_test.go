package cmd_test

import (
	"github.com/oslokommune/aws-delete-stacks/delete_stacks"
	"testing"
)

func TestFilter(t *testing.T) {
	testCases := []*DeleteCmdTestCase{
		{
			name: "Should not delete stacks that doesn't contain include filter",
			args: "myenv",
			stackSummaries: []*delete_stacks.StackSummary{
				newStackSummary("somestack-dev", mockConstants().StackStatusCreateComplete),
				newStackSummary("some-other-stack-dev", mockConstants().StackStatusCreateComplete),
			},
			expectedStdoutContains: []string{"Would delete 0 stack(s)"},
		},
		{
			name: "Should delete stacks that contains include filter",
			args: "other-stack",
			stackSummaries: []*delete_stacks.StackSummary{
				newStackSummary("somestack-dev", mockConstants().StackStatusCreateComplete),
				newStackSummary("some-other-stack-dev", mockConstants().StackStatusCreateComplete),
			},
			expectedStdoutContains: []string{
				"Would delete 1 stack(s)",
				"some-other-stack-dev",
			},
		},
		{
			name: "Should not delete stacks that contains include filter and exclude filter",
			args: "other-stack --exclude hosted-zone",
			stackSummaries: []*delete_stacks.StackSummary{
				newStackSummary("somestack-dev", mockConstants().StackStatusCreateComplete),
				newStackSummary("some-other-stack-dev", mockConstants().StackStatusCreateComplete),
				newStackSummary("some-other-stack-hosted-zone-dev", mockConstants().StackStatusCreateComplete),
			},
			expectedStdoutContains: []string{
				"Would delete 1 stack(s)",
				"some-other-stack-dev",
			},
		},
	}

	runTestCases(t, testCases)
}