package cmd_test

import (
	"github.com/oslokommune/aws-delete-stacks/cloudformation_api"
	"time"
)

func NewCloudFormationMock() *CloudFormationMock {
	return &CloudFormationMock{}
}

type CloudFormationMock struct {
	StackSummaries []*cloudformation_api.StackSummary
}

func (c *CloudFormationMock) Constants() *cloudformation_api.Constants {
	return mockConstants()
}

func mockConstants() *cloudformation_api.Constants {
	return &cloudformation_api.Constants{
		StackStatusDeleteInProgress: "IN_PROGRESS",
		StackStatusCreateComplete:   "CREATE_COMPLETE",
		StackStatusDeleteFailed:     "DELETE_FAILED",
		StackStatusDeleteComplete:   "DELETE_COMPLETE",
	}
}

func (c *CloudFormationMock) ListStacks(_ *cloudformation_api.ListStacksInput) (*cloudformation_api.ListStacksOutput, error) {
	return &cloudformation_api.ListStacksOutput{
		StackSummaries: c.StackSummaries,
		NextToken:      nil,
	}, nil
}

func newStackSummary(stackName string, stackStatus string) *cloudformation_api.StackSummary {
	id := "id-" + stackName
	loc := time.Local
	t := time.Date(2020, 1, 1, 0, 0, 0, 0, loc)

	return &cloudformation_api.StackSummary{
		StackName:    &stackName,
		StackId:      &id,
		StackStatus:  &stackStatus,
		CreationTime: &t,
	}
}

func (c *CloudFormationMock) DeleteStack(_ *cloudformation_api.DeleteStackInput) (*cloudformation_api.DeleteStackOutput, error) {
	panic("implement me")
}

func (c *CloudFormationMock) DescribeStacks(_ *cloudformation_api.DescribeStacksInput) (*cloudformation_api.DescribeStacksOutput, error) {
	panic("implement me")
}
