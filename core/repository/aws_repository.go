package repository

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

type CloudFormationRepository interface {
	ListStacks(input *ListStacksInput) (*ListStacksOutput, error)
	DeleteStack(input *DeleteStackInput) (*DeleteStackOutput, error)
	DescribeStacks(input *DescribeStacksInput) (*DescribeStacksOutput, error)
}

type AWSCloudFormationRepository struct {
	cloudformation *cloudformation.CloudFormation
}

func NewAWSCloudFormationRepository(sess *session.Session) *AWSCloudFormationRepository {
	cf := cloudformation.New(sess)

	return &AWSCloudFormationRepository{
		cloudformation: cf,
	}
}

func (r *AWSCloudFormationRepository) ListStacks(input *ListStacksInput) (*ListStacksOutput, error) {
	cfInput := &cloudformation.ListStacksInput{
		NextToken:         input.NextToken,
		StackStatusFilter: input.StackStatusFilter,
	}

	cfOutput, err := r.cloudformation.ListStacks(cfInput)
	if err != nil {
		return nil, err
	}

	output := &ListStacksOutput{}
	output.NextToken = cfOutput.NextToken
	output.StackSummaries = make([]*StackSummary, len(cfOutput.StackSummaries))

	for i, summary := range cfOutput.StackSummaries {
		output.StackSummaries[i] = &StackSummary{
			StackName:    summary.StackName,
			StackId:      summary.StackId,
			StackStatus:  summary.StackStatus,
			CreationTime: summary.CreationTime,
		}
	}

	return output, nil
}

func (r *AWSCloudFormationRepository) DeleteStack(input *DeleteStackInput) (*DeleteStackOutput, error) {
	cfInput := &cloudformation.DeleteStackInput{
		ClientRequestToken: input.ClientRequestToken,
		StackName:          input.StackName,
	}

	_, err := r.cloudformation.DeleteStack(cfInput)
	if err != nil {
		return nil, err
	}

	return &DeleteStackOutput{}, nil
}

func (r *AWSCloudFormationRepository) DescribeStacks(input *DescribeStacksInput) (*DescribeStacksOutput, error) {
	cfInput := &cloudformation.DescribeStacksInput{
		NextToken: input.NextToken,
		StackName: input.StackName,
	}

	cfOutput, err := r.cloudformation.DescribeStacks(cfInput)
	if err != nil {
		return nil, err
	}

	output := &DescribeStacksOutput{}
	output.Stacks = make([]*Stack, len(cfOutput.Stacks))

	for i, stack := range cfOutput.Stacks {
		output.Stacks[i] = &Stack{
			StackStatus: stack.StackStatus,
		}
	}

	return output, nil
}
