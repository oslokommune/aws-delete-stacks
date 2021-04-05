package cmd_test

import (
	"github.com/oslokommune/aws-delete-stacks/core/repository"
	"testing"
)

func TestFilter(t *testing.T) {
	testCases := []*DeleteCmdTestCase{
		{
			Name: "Should not delete stacks that doesn't contain include filter",
			Args: "myenv",
			ListStacksOutput: BuildListStackOutput(nil,
				NewStackSummary("somestack-dev", repository.StackStatusCreateComplete),
				NewStackSummary("some-other-stack-dev", repository.StackStatusCreateComplete),
			),
			ExpectedStdoutContains: []string{"Would delete 0 stack(s)"},
		},
		{
			Name: "Should delete stacks that contains include filter",
			Args: "mystack",
			ListStacksOutput: BuildListStackOutput(nil,
				NewStackSummary("some-first-stack", repository.StackStatusCreateComplete),
				NewStackSummary("mystack", repository.StackStatusCreateComplete),
				NewStackSummary("some-other-stack-dev", repository.StackStatusCreateComplete),
			),
			ExpectedStdoutContains: []string{
				"Would delete 1 stack(s)",
				"mystack",
			},
		},
		{
			Name: "Should not delete stacks that contains include filter and exclude filter",
			Args: "other-stack --exclude hosted-zone",
			ListStacksOutput: BuildListStackOutput(nil,
				NewStackSummary("firststack", repository.StackStatusCreateComplete),
				NewStackSummary("some-other-stack-dev", repository.StackStatusCreateComplete),
				NewStackSummary("some-other-stack-hosted-zone-dev", repository.StackStatusCreateComplete),
				NewStackSummary("laststack", repository.StackStatusCreateComplete),
			),
			ExpectedStdoutContains: []string{
				"Would delete 1 stack(s)",
				"some-other-stack-dev",
			},
		},
	}

	RunTestCases(t, testCases)
}
