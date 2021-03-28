package aws_cloudformation_old

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func NewAWSCloudFormationOld() *AWSCloudFormationOld {
	t := true
	s := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			CredentialsChainVerboseErrors: &t,
		},
		SharedConfigState: session.SharedConfigEnable,
	}))

	cf := cloudformation.New(s)

	return &AWSCloudFormationOld{
		cloudformation: cf,
	}
}

type AWSCloudFormationOld struct {
	cloudformation *cloudformation.CloudFormation
}

func (c *AWSCloudFormationOld) ListStacks(input *ListStacksInput) (*ListStacksOutput, error) {
	outputOrg, err := c.cloudformation.ListStacks(&cloudformation.ListStacksInput{
		NextToken:         input.NextToken,
		StackStatusFilter: input.StackStatusFilter,
	})
	if err != nil {
		return nil, err
	}

	output := &ListStacksOutput{}
	output.NextToken = outputOrg.NextToken
	output.StackSummaries = make([]*StackSummary, len(outputOrg.StackSummaries))

	for i, summary := range outputOrg.StackSummaries {
		output.StackSummaries[i] = &StackSummary{
			StackName:    summary.StackName,
			StackId:      summary.StackId,
			StackStatus:  summary.StackStatus,
			CreationTime: summary.CreationTime,
		}
	}

	return output, nil
}

func (c *AWSCloudFormationOld) DeleteStack(input *DeleteStackInput) (*DeleteStackOutput, error) {
	_, err := c.cloudformation.DeleteStack(&cloudformation.DeleteStackInput{
		ClientRequestToken: input.ClientRequestToken,
		StackName:          input.StackName,
	})
	if err != nil {
		return nil, err
	}

	return &DeleteStackOutput{}, nil
}

func (c *AWSCloudFormationOld) DescribeStacks(input *DescribeStacksInput) (*DescribeStacksOutput, error) {
	outputOrg, err := c.cloudformation.DescribeStacks(&cloudformation.DescribeStacksInput{
		NextToken: input.NextToken,
		StackName: input.StackName,
	})
	if err != nil {
		return nil, err
	}

	output := &DescribeStacksOutput{}
	output.Stacks = make([]*Stack, len(outputOrg.Stacks))

	for i, stack := range outputOrg.Stacks {
		output.Stacks[i] = &Stack{
			StackStatus: stack.StackStatus,
		}
	}

	return output, nil
}
