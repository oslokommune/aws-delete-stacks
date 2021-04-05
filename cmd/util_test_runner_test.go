package cmd_test

import (
	"bytes"
	"github.com/oslokommune/aws-delete-stacks/cmd"
	"github.com/oslokommune/aws-delete-stacks/core/application"
	"github.com/oslokommune/aws-delete-stacks/core/repository"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

type DeleteCmdTestCase struct {
	Name                   string
	Args                   string
	ListStacksOutput       *repository.ListStacksOutput
	ExpectError            bool
	ExpectedStdoutContains []string
	ExpectedStderrContains []string
}

func RunTestCases(t *testing.T, testCases []*DeleteCmdTestCase) {
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			// Given
			mock := NewCloudFormationRepositoryMock()
			mock.ListStackOutput = tc.ListStacksOutput

			// When
			stdout, stderr := runDeleteCommand(t, mock, tc.Args)

			// Then
			if tc.ExpectError {
				for _, ec := range tc.ExpectedStderrContains {
					assert.Contains(t, stderr, ec)
				}
			} else {
				for _, ec := range tc.ExpectedStdoutContains {
					assert.Contains(t, stdout, ec)
				}
			}
		})
	}
}

func BuildListStackOutput(nextToken *string, summaries ...*repository.StackSummary) *repository.ListStacksOutput {
	return &repository.ListStacksOutput{
		StackSummaries: summaries,
		NextToken:      nextToken,
	}
}

func runDeleteCommand(t *testing.T, repository repository.CloudFormationRepository, argsString string) (string, string) {
	cloudFormation := application.NewCloudFormation(repository)

	deleteAndArgs := strings.Trim(argsString, " ")
	args := strings.Split(deleteAndArgs, " ")

	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer

	c := cmd.BuildDeleteCommand(cloudFormation)
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
