package cmd_test

import (
	"github.com/oslokommune/aws-delete-stacks/core/delete_stacks/cloudformation_api"
	"time"
)

type CloudFormationMock struct {
	stacks []*cloudformation_api.Stack
}

func (c *CloudFormationMock) ListStacks(_ []cloudformation_api.StackStatus) ([]*cloudformation_api.Stack, error) {
	return c.stacks, nil
}

func (c *CloudFormationMock) DeleteStack(_ *cloudformation_api.Stack) error {
	panic("implement me")
}

func (c *CloudFormationMock) GetStackStatus(_ *cloudformation_api.Stack) (cloudformation_api.StackStatus, error) {
	panic("implement me")
}

func NewCloudFormationMock() *CloudFormationMock {
	return &CloudFormationMock{}
}

func newStack(stackName string, stackStatus cloudformation_api.StackStatus) *cloudformation_api.Stack {
	loc := time.Local
	t := time.Date(2020, 1, 1, 0, 0, 0, 0, loc)

	return &cloudformation_api.Stack{
		Name:         stackName,
		Status:       stackStatus,
		CreationTime: &t,
	}
}
