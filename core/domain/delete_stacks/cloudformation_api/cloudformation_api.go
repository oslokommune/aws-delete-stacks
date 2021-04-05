package cloudformation_api

type CloudFormation interface {
	ListStacks(statusFilter []StackStatus) ([]*Stack, error)
	DeleteStack(input *Stack) error
	GetStackStatus(input *Stack) (StackStatus, error)
}
