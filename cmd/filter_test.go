package cmd_test

import (
	"github.com/oslokommune/aws-delete-stacks/core/delete_stacks/cloudformation_api"
	"testing"
)

func TestFilter(t *testing.T) {
	testCases := []*DeleteCmdTestCase{
		{
			name: "Should not delete stacks that doesn't contain include filter",
			args: "myenv",
			stacks: []*cloudformation_api.Stack{
				newStack("somestack-dev", cloudformation_api.StackStatusCreateComplete),
				newStack("some-other-stack-dev", cloudformation_api.StackStatusCreateComplete),
			},
			expectedStdoutContains: []string{"Would delete 0 stack(s)"},
		},
		{
			name: "Should delete stacks that contains include filter",
			args: "mystack",
			stacks: []*cloudformation_api.Stack{
				newStack("some-first-stack", cloudformation_api.StackStatusCreateComplete),
				newStack("mystack", cloudformation_api.StackStatusCreateComplete),
				newStack("some-other-stack-dev", cloudformation_api.StackStatusCreateComplete),
			},
			expectedStdoutContains: []string{
				"Would delete 1 stack(s)",
				"mystack",
			},
		},
		{
			name: "Should not delete stacks that contains include filter and exclude filter",
			args: "other-stack --exclude hosted-zone",
			stacks: []*cloudformation_api.Stack{
				newStack("firststack", cloudformation_api.StackStatusCreateComplete),
				newStack("some-other-stack-dev", cloudformation_api.StackStatusCreateComplete),
				newStack("some-other-stack-hosted-zone-dev", cloudformation_api.StackStatusCreateComplete),
				newStack("laststack", cloudformation_api.StackStatusCreateComplete),
			},
			expectedStdoutContains: []string{
				"Would delete 1 stack(s)",
				"some-other-stack-dev",
			},
		},
	}

	runTestCases(t, testCases)
}
