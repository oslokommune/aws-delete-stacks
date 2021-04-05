package cmd_test

import (
	"github.com/oslokommune/aws-delete-stacks/core/repository"
	"time"
)

type CloudFormationRepositoryMock struct {
	ListStackOutput *repository.ListStacksOutput
}

func NewCloudFormationRepositoryMock() *CloudFormationRepositoryMock {
	return &CloudFormationRepositoryMock{}
}

func (m *CloudFormationRepositoryMock) ListStacks(_ *repository.ListStacksInput) (*repository.ListStacksOutput, error) {
	return m.ListStackOutput, nil
}

func (m *CloudFormationRepositoryMock) DeleteStack(_ *repository.DeleteStackInput) (*repository.DeleteStackOutput, error) {
	panic("implement me")
}

func (m *CloudFormationRepositoryMock) DescribeStacks(_ *repository.DescribeStacksInput) (*repository.DescribeStacksOutput, error) {
	panic("implement me")
}

func NewStackSummary(stackName string, stackStatus string) *repository.StackSummary {
	id := "id-" + stackName
	loc := time.Local
	t := time.Date(2020, 1, 1, 0, 0, 0, 0, loc)

	return &repository.StackSummary{
		StackName:    &stackName,
		StackId:      &id,
		StackStatus:  &stackStatus,
		CreationTime: &t,
	}
}
