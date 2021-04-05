package repository

import (
	"fmt"
	"time"
)

const (
	StackStatusCreateComplete   = "CREATE_COMPLETE"
	StackStatusDeleteInProgress = "DELETE_IN_PROGRESS"
	StackStatusDeleteComplete   = "DELETE_COMPLETE"
	StackStatusDeleteFailed     = "DELETE_FAILED"
)

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

type RequestFailure interface {
	Error

	// The status code of the HTTP response.
	StatusCode() int

	// The request ID returned by the service for a request failure. This will
	// be empty if no request ID is available such as the request failed due
	// to a connection error.
	RequestID() string
}

type Error interface {
	// Satisfy the generic error interface.
	error

	// Returns the short phrase depicting the classification of the error.
	Code() string

	// Returns the error details message.
	Message() string

	// Returns the original error if one was set.  Nil is returned if not set.
	OrigErr() error
}
