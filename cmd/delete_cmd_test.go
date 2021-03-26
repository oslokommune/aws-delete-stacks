package cmd_test

import (
	"bytes"
	"github.com/oslokommune/aws-delete-stacks/cmd"
	"github.com/oslokommune/aws-delete-stacks/delete_stacks"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

type DeleteCmdTestCase struct {
	name                   string
	args                   string
	stackSummaries         []*delete_stacks.StackSummary
	expectError            bool
	expectedStdoutContains []string
	expectedStderrContains []string
}

func TestDeleteCmd(t *testing.T) {
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

func runTestCases(t *testing.T, testCases []*DeleteCmdTestCase) {
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Given
			mock := NewCloudFormationMock()
			mock.StackSummaries = tc.stackSummaries

			// When
			stdout, stderr := runDeleteCommand(t, mock, tc.args)

			// Then
			if tc.expectError {
				for _, ec := range tc.expectedStderrContains {
					assert.Contains(t, stderr, ec)
				}
			} else {
				for _, ec := range tc.expectedStdoutContains {
					assert.Contains(t, stdout, ec)
				}
			}
		})
	}
}

func runDeleteCommand(t *testing.T, cf delete_stacks.CloudFormation, argsString string) (string, string) {
	deleteAndArgs := strings.Trim(argsString, " ")
	args := strings.Split(deleteAndArgs, " ")

	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer

	c := cmd.BuildDeleteCommand(cf)
	c.SetArgs(args)
	c.SetOut(&stdoutBuffer)
	c.SetErr(&stderrBuffer)

	err := c.Execute()

	stdout, err := ioutil.ReadAll(&stdoutBuffer)
	assert.NoError(t, err)

	stderr, err := ioutil.ReadAll(&stderrBuffer)
	assert.NoError(t, err)

	return string(stdout), string(stderr)
}
