package application

import (
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/oslokommune/aws-delete-stacks/core/domain/delete_stacks/cloudformation_api"
	"github.com/oslokommune/aws-delete-stacks/core/repository"
	"strings"
)

func NewCloudFormation(repository repository.CloudFormationRepository) cloudformation_api.CloudFormation {
	return &CloudFormation{
		repository: repository,
	}
}

type CloudFormation struct {
	repository repository.CloudFormationRepository
}

func (c *CloudFormation) ListStacks(statusFilter []cloudformation_api.StackStatus) ([]*cloudformation_api.Stack, error) {
	apiStatusFilter := c.toApiStatusFilter(statusFilter)
	outputs := make([]*repository.StackSummary, 0)

	i := 0
	crashProtection := 1000
	var nextPageToken *string

	for i < crashProtection {
		var output *repository.ListStacksOutput
		var err error

		if i == 0 {
			output, err = c.repository.ListStacks(&repository.ListStacksInput{
				StackStatusFilter: apiStatusFilter,
			})
		} else {
			output, err = c.repository.ListStacks(&repository.ListStacksInput{
				NextToken:         nextPageToken,
				StackStatusFilter: apiStatusFilter,
			})
		}

		if err != nil {
			return nil, fmt.Errorf("cloudformation list stacks: %w", err)
		}

		outputs = append(outputs, output.StackSummaries...)

		if output.NextToken == nil {
			break
		}

		nextPageToken = output.NextToken
		i++
	}

	return toStacks(outputs)
}

func (c *CloudFormation) toApiStatusFilter(statusesDomain []cloudformation_api.StackStatus) []*string {
	statusesAws := make([]*string, len(statusesDomain))

	for i, status := range statusesDomain {
		s := status.String()
		statusesAws[i] = &s
	}

	return statusesAws
}

func toStacks(summaries []*repository.StackSummary) ([]*cloudformation_api.Stack, error) {
	stacks := make([]*cloudformation_api.Stack, len(summaries))

	for i, summary := range summaries {
		stackStatus, err := cloudformation_api.ParseStackStatus(*summary.StackStatus)
		if err != nil {
			return nil, fmt.Errorf("parse stack status: %w", err)
		}

		stacks[i] = &cloudformation_api.Stack{
			Name:         *summary.StackName,
			Status:       stackStatus,
			CreationTime: summary.CreationTime,
		}
	}

	return stacks, nil
}

func (c *CloudFormation) DeleteStack(stack *cloudformation_api.Stack) error {
	clientRequestToken, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("generating uuid: %w", err)
	}

	_, err = c.repository.DeleteStack(&repository.DeleteStackInput{
		ClientRequestToken: &clientRequestToken,
		StackName:          &stack.Name,
	})
	if err != nil {
		return fmt.Errorf("cloudformation delete stack: %w", err)
	}

	return nil
}

func (c *CloudFormation) GetStackStatus(stack *cloudformation_api.Stack) (cloudformation_api.StackStatus, error) {
	describeStackOutputs, err := c.repository.DescribeStacks(&repository.DescribeStacksInput{
		StackName: &stack.Name,
	})
	if err != nil {
		awsError, ok := err.(repository.RequestFailure)
		if !ok {
			return "", fmt.Errorf("cloudformation describe stack: %w", err)
		}

		if strings.Contains(awsError.Message(), "does not exist") {
			return cloudformation_api.StackStatusDeleteComplete, nil
		}

		return "", fmt.Errorf("cloudformation describe stack: %w", err)
	}

	if len(describeStackOutputs.Stacks) > 1 {
		names := c.getStackNames(describeStackOutputs)
		return "", fmt.Errorf(
			"internal error, expected 1 stack, but got %d: %s", len(describeStackOutputs.Stacks), names)
	}

	describeStackOutput := describeStackOutputs.Stacks[0]

	return toStackStatus(*describeStackOutput.StackStatus)
}

func toStackStatus(s string) (cloudformation_api.StackStatus, error) {
	switch s {
	case repository.StackStatusCreateComplete:
		return cloudformation_api.StackStatusCreateComplete, nil
	case repository.StackStatusDeleteInProgress:
		return cloudformation_api.StackStatusDeleteInProgress, nil
	case repository.StackStatusDeleteComplete:
		return cloudformation_api.StackStatusDeleteComplete, nil
	case repository.StackStatusDeleteFailed:
		return cloudformation_api.StackStatusDeleteFailed, nil
	}

	return "", fmt.Errorf("could not find correct enum for string '%s'", s)
}

func (c *CloudFormation) getStackNames(describeStackOutput *repository.DescribeStacksOutput) string {
	nameSlice := make([]string, len(describeStackOutput.Stacks))
	for i, s := range describeStackOutput.Stacks {
		nameSlice[i] = *s.StackName
	}

	names := strings.Join(nameSlice, ", ")
	return names
}
