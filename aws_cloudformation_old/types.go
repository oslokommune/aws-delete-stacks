package aws_cloudformation_old

import (
	"fmt"
	"time"
)

type CloudFormation interface {
	Constants() *Constants

	ListStacks(input *ListStacksInput) (*ListStacksOutput, error)
	DeleteStack(input *DeleteStackInput) (*DeleteStackOutput, error)
	DescribeStacks(input *DescribeStacksInput) (*DescribeStacksOutput, error)
}

type Constants struct {
	StackStatusDeleteInProgress string
	StackStatusCreateComplete   string
	StackStatusDeleteFailed     string
	StackStatusDeleteComplete   string
}

type ListStacksInput struct {
	StackStatusFilter []*string
	NextToken         *string
}

type ListStacksOutput struct {
	StackSummaries []*StackSummary
	NextToken      *string
}

type StackSummary struct {
	StackName    *string
	StackId      *string
	StackStatus  *string
	CreationTime *time.Time
}

type DeleteStackInput struct {
	ClientRequestToken *string
	StackName          *string
}

type DeleteStackOutput struct{}

type DescribeStacksInput struct {
	StackName *string
	NextToken *string
}

type DescribeStacksOutput struct {
	Stacks []*Stack
}

type Stack struct {
	StackName   *string
	StackStatus *string
	RoleARN     *string
	Time        *time.Time
}

func (i *Stack) String() string {
	return fmt.Sprintf("%s (%s)", *i.StackName, i.Time.Format(time.RFC822))
}
