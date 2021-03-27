package cloudformation_api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/oslokommune/aws-delete-stacks/delete_stacks"
)

func NewAWSCloudFormation() *AWSCloudFormation {
	t := true
	s := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			CredentialsChainVerboseErrors: &t,
		},
		SharedConfigState: session.SharedConfigEnable,
	}))

	cf := cloudformation.New(s)

	return &AWSCloudFormation{
		cloudformation: cf,
	}
}

type AWSCloudFormation struct {
	cloudformation *cloudformation.CloudFormation
}

func constants() *delete_stacks.Constants {
	return &delete_stacks.Constants{
		StackStatusDeleteInProgress: cloudformation.StackStatusDeleteInProgress,
		StackStatusCreateComplete:   cloudformation.StackStatusCreateComplete,
		StackStatusDeleteFailed:     cloudformation.StackStatusDeleteFailed,
		StackStatusDeleteComplete:   cloudformation.StackStatusDeleteComplete,
	}
}

func (c *AWSCloudFormation) Constants() *delete_stacks.Constants {
	return constants()
}

func (c *AWSCloudFormation) ListStacks(input *delete_stacks.ListStacksInput) (*delete_stacks.ListStacksOutput, error) {
	outputOrg, err := c.cloudformation.ListStacks(&cloudformation.ListStacksInput{
		NextToken:         input.NextToken,
		StackStatusFilter: input.StackStatusFilter,
	})
	if err != nil {
		return nil, err
	}

	output := &delete_stacks.ListStacksOutput{}
	output.NextToken = outputOrg.NextToken
	output.StackSummaries = make([]*delete_stacks.StackSummary, len(outputOrg.StackSummaries))

	for i, summary := range outputOrg.StackSummaries {
		output.StackSummaries[i] = &delete_stacks.StackSummary{
			StackName:    summary.StackName,
			StackId:      summary.StackId,
			StackStatus:  summary.StackStatus,
			CreationTime: summary.CreationTime,
		}
	}

	return output, nil
}

func (c *AWSCloudFormation) DeleteStack(input *delete_stacks.DeleteStackInput) (*delete_stacks.DeleteStackOutput, error) {
	_, err := c.cloudformation.DeleteStack(&cloudformation.DeleteStackInput{
		ClientRequestToken: input.ClientRequestToken,
		StackName:          input.StackName,
	})
	if err != nil {
		return nil, err
	}

	return &delete_stacks.DeleteStackOutput{}, nil
}

func (c *AWSCloudFormation) DescribeStacks(input *delete_stacks.DescribeStacksInput) (*delete_stacks.DescribeStacksOutput, error) {
	outputOrg, err := c.cloudformation.DescribeStacks(&cloudformation.DescribeStacksInput{
		NextToken: input.NextToken,
		StackName: input.StackName,
	})
	if err != nil {
		return nil, err
	}

	output := &delete_stacks.DescribeStacksOutput{}
	output.Stacks = make([]*delete_stacks.Stack, len(outputOrg.Stacks))

	for i, stack := range outputOrg.Stacks {
		output.Stacks[i] = &delete_stacks.Stack{
			StackStatus: stack.StackStatus,
		}
	}

	return output, nil
}
