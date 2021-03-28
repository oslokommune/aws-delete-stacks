package cmd_test

import (
	"github.com/oslokommune/aws-delete-stacks/aws_cloudformation_old"
	"time"
)

type CloudFormationMockOld struct {
	StackSummaries []*aws_cloudformation_old.StackSummary
}

func (c *CloudFormationMockOld) ListStacks(_ *aws_cloudformation_old.ListStacksInput) (*aws_cloudformation_old.ListStacksOutput, error) {
	return &aws_cloudformation_old.ListStacksOutput{
		StackSummaries: c.StackSummaries,
		NextToken:      nil,
	}, nil
}

func newStackSummary(stackName string, stackStatus string) *aws_cloudformation_old.StackSummary {
	id := "id-" + stackName
	loc := time.Local
	t := time.Date(2020, 1, 1, 0, 0, 0, 0, loc)

	return &aws_cloudformation_old.StackSummary{
		StackName:    &stackName,
		StackId:      &id,
		StackStatus:  &stackStatus,
		CreationTime: &t,
	}
}

func (c *CloudFormationMockOld) DeleteStack(_ *aws_cloudformation_old.DeleteStackInput) (*aws_cloudformation_old.DeleteStackOutput, error) {
	panic("implement me")
}

func (c *CloudFormationMockOld) DescribeStacks(_ *aws_cloudformation_old.DescribeStacksInput) (*aws_cloudformation_old.DescribeStacksOutput, error) {
	panic("implement me")
}
