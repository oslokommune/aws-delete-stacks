package cmd_test

import (
	"bytes"
	"fmt"
	"github.com/oslokommune/aws-delete-stacks/cmd"
	"github.com/oslokommune/aws-delete-stacks/delete_stacks"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

func TestDeleteCmd(t *testing.T) {
	testCases := []struct {
		name                   string
		includeFilter          string
		stackSummaries         []*delete_stacks.StackSummary
		expectedOutputContains string
	}{
		{
			name:          "Should not delete stacks that doesn't contain include filter",
			includeFilter: "myenv",
			stackSummaries: []*delete_stacks.StackSummary{
				newStackSummary("somestack-dev", mockConstants().StackStatusCreateComplete),
				newStackSummary("some-other-stack-dev", mockConstants().StackStatusCreateComplete),
			},
			expectedOutputContains: "Would delete 0 stack(s)",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Given
			mock := NewCloudFormationMock()
			mock.StackSummaries = tc.stackSummaries

			// When
			//runDeleteCommand(t,"--include mystack --exclude mystack-importantresource")
			cmdOutput := runDeleteCommand(t, mock, tc.includeFilter)

			// Then
			assert.Contains(t, cmdOutput, tc.expectedOutputContains)
			fmt.Println(cmdOutput)
		})
	}
}

func runDeleteCommand(t *testing.T, cf delete_stacks.CloudFormation, argsString string) string {
	args := strings.Split(argsString, " ")

	var buf bytes.Buffer

	c := cmd.BuildDeleteCommand(cf)

	c.SetArgs(args)
	c.SetOut(&buf)

	err := c.Execute()
	assert.NoError(t, err)

	out, err := ioutil.ReadAll(&buf)
	assert.NoError(t, err)

	return string(out)
}
