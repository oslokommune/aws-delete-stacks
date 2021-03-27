package cmd_test

import (
	"bytes"
	"github.com/oslokommune/aws-delete-stacks/cloudformation_api"
	"github.com/oslokommune/aws-delete-stacks/cmd"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

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

func runDeleteCommand(t *testing.T, cf cloudformation_api.CloudFormation, argsString string) (string, string) {
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
