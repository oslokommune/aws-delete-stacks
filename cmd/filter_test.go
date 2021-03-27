package cmd_test

import (
	"github.com/oslokommune/aws-delete-stacks/cloudformation_api"
	"testing"
)

func TestFilter(t *testing.T) {
	testCases := []*DeleteCmdTestCase{
		{
			name: "Should not delete stacks that doesn't contain include filter",
			args: "myenv",
			stackSummaries: []*cloudformation_api.StackSummary{
				newStackSummary("somestack-dev", mockConstants().StackStatusCreateComplete),
				newStackSummary("some-other-stack-dev", mockConstants().StackStatusCreateComplete),
			},
			expectedStdoutContains: []string{"Would delete 0 stack(s)"},
		},
		{
			name: "Should delete stacks that contains include filter",
			args: "mystack",
			stackSummaries: []*cloudformation_api.StackSummary{
				newStackSummary("some-first-stack", mockConstants().StackStatusCreateComplete),
				newStackSummary("mystack", mockConstants().StackStatusCreateComplete),
				newStackSummary("some-other-stack-dev", mockConstants().StackStatusCreateComplete),
			},
			expectedStdoutContains: []string{
				"Would delete 1 stack(s)",
				"mystack",
			},
		},
		{
			name: "Should not delete stacks that contains include filter and exclude filter",
			args: "other-stack --exclude hosted-zone",
			stackSummaries: []*cloudformation_api.StackSummary{
				newStackSummary("firststack", mockConstants().StackStatusCreateComplete),
				newStackSummary("some-other-stack-dev", mockConstants().StackStatusCreateComplete),
				newStackSummary("some-other-stack-hosted-zone-dev", mockConstants().StackStatusCreateComplete),
				newStackSummary("laststack", mockConstants().StackStatusCreateComplete),
			},
			expectedStdoutContains: []string{
				"Would delete 1 stack(s)",
				"some-other-stack-dev",
			},
		},
	}

	runTestCases(t, testCases)
}
