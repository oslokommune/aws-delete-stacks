package delete_stacks

import "time"

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
	StackStatus *string
}
