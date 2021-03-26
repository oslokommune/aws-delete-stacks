package cmd_test

import (
	"github.com/oslokommune/aws-delete-stacks/delete_stacks"
	"time"
)

func NewCloudFormationMock() *CloudFormationMock {
	return &CloudFormationMock{}
}

type CloudFormationMock struct {
	StackSummaries []*delete_stacks.StackSummary
}

func mockConstants() *delete_stacks.Constants {
	return &delete_stacks.Constants{
		StackStatusDeleteInProgress: "IN_PROGRESS",
		StackStatusCreateComplete:   "CREATE_COMPLETE",
		StackStatusDeleteFailed:     "DELETE_FAILED",
	}
}

func (c *CloudFormationMock) Constants() *delete_stacks.Constants {
	return mockConstants()
}

func (c *CloudFormationMock) ListStacks(_ *delete_stacks.ListStacksInput) (*delete_stacks.ListStacksOutput, error) {
	return &delete_stacks.ListStacksOutput{
		StackSummaries: c.StackSummaries,
		NextToken:      nil,
	}, nil
}

func newStackSummary(stackName string, stackStatus string) *delete_stacks.StackSummary {
	id := "id-" + stackName
	loc := time.Local
	t := time.Date(2020, 1, 1, 0, 0, 0, 0, loc)

	return &delete_stacks.StackSummary{
		StackName:    &stackName,
		StackId:      &id,
		StackStatus:  &stackStatus,
		CreationTime: &t,
	}
}

func (c *CloudFormationMock) DeleteStack(_ *delete_stacks.DeleteStackInput) (*delete_stacks.DeleteStackOutput, error) {
	panic("implement me")
}

func (c *CloudFormationMock) DescribeStacks(_ *delete_stacks.DescribeStacksInput) (*delete_stacks.DescribeStacksOutput, error) {
	panic("implement me")
}
