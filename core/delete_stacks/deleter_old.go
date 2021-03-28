package delete_stacks

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/go-uuid"
	"github.com/oslokommune/aws-delete-stacks/aws_cloudformation_old"
	"io"
	"strings"
	"time"
)

type DeleterOld struct {
	cf  aws_cloudformation_old.CloudFormation
	out io.Writer
}

func NewDeleterOld(cf aws_cloudformation_old.CloudFormation, out io.Writer) *DeleterOld {
	return &DeleterOld{
		cf:  cf,
		out: out,
	}
}

func (d *DeleterOld) DeleteCloudFormationStacks(includeFilter string, excludeFilter string, force bool) error {
	output, err := d.listStacks()
	if err != nil {
		return err
	}

	var toDelete []*aws_cloudformation_old.Stack
	toDelete = d.toStacks(output)

	toDelete, err = d.filter(toDelete, includeFilter, excludeFilter)
	if err != nil {
		return fmt.Errorf("filter stack input: %w", err)
	}

	err = d.deleteStacks(toDelete, force)
	if err != nil {
		return fmt.Errorf("delete stacks: %w", err)
	}

	return nil
}

func (d *DeleterOld) listStacks() ([]*aws_cloudformation_old.StackSummary, error) {
	outputs := make([]*aws_cloudformation_old.StackSummary, 0)

	statusDeleteInProgress := d.cf.Constants().StackStatusDeleteInProgress
	statusCreateComplete := d.cf.Constants().StackStatusCreateComplete
	statusDeleteFailed := d.cf.Constants().StackStatusDeleteFailed
	statusFilter := []*string{&statusDeleteInProgress, &statusCreateComplete, &statusDeleteFailed}

	i := 0
	crashProtection := 1000
	var nextPageToken *string

	for i < crashProtection {
		var output *aws_cloudformation_old.ListStacksOutput
		var err error

		if i == 0 {
			output, err = d.cf.ListStacks(&aws_cloudformation_old.ListStacksInput{
				StackStatusFilter: statusFilter,
			})
		} else {
			output, err = d.cf.ListStacks(&aws_cloudformation_old.ListStacksInput{
				NextToken:         nextPageToken,
				StackStatusFilter: statusFilter,
			})
		}

		if err != nil {
			return nil, fmt.Errorf("list stacks: %w", err)
		}

		outputs = append(outputs, output.StackSummaries...)

		if output.NextToken == nil {
			break
		}
		nextPageToken = output.NextToken
		i++
	}

	return outputs, nil
}

func (d *DeleterOld) toStacks(output []*aws_cloudformation_old.StackSummary) []*aws_cloudformation_old.Stack {
	input := make([]*aws_cloudformation_old.Stack, 0)

	for _, summary := range output {
		input = append(input, &aws_cloudformation_old.Stack{
			StackName:   summary.StackName,
			RoleARN:     summary.StackId,
			StackStatus: summary.StackStatus,
			Time:        summary.CreationTime,
		})
	}

	return input
}

func (d *DeleterOld) filter(stacks []*aws_cloudformation_old.Stack, includeFilter string, excludeFilter string) ([]*aws_cloudformation_old.Stack, error) {
	filtered := make([]*aws_cloudformation_old.Stack, 0)

	var add *aws_cloudformation_old.Stack
	for _, input := range stacks {
		add = nil

		if len(includeFilter) > 0 && strings.Contains(*input.StackName, includeFilter) {
			add = input
		}

		if len(excludeFilter) > 0 && strings.Contains(*input.StackName, excludeFilter) {
			add = nil
		}

		if add != nil {
			filtered = append(filtered, input)
		}
	}

	return filtered, nil
}

func (d *DeleterOld) deleteStacks(stacks []*aws_cloudformation_old.Stack, force bool) error {
	if force {
		_, err := fmt.Fprintf(d.out, "- Deleting %d stacks\n", len(stacks))
		if err != nil {
			return err
		}
	} else {
		_, err := fmt.Fprintf(d.out, "- Would delete %d stack(s)\n", len(stacks))
		if err != nil {
			return err
		}
	}

	for _, stack := range stacks {
		if !force {
			_, err := fmt.Fprintln(d.out, stack)
			if err != nil {
				return err
			}

			continue
		}

		err := d.deleteStack(stack)
		if err != nil {
			return err
		}
	}

	_, err := fmt.Fprintf(d.out, "\nNo CloudFormation stacks were deleted, as you didn't specify the --force flag.\n")
	if err != nil {
		return err
	}

	return nil
}

func (d *DeleterOld) deleteStack(stack *aws_cloudformation_old.Stack) error {
	_, err := fmt.Fprintf(d.out, "\nDeleting stack: %+v\n", stack)
	if err != nil {
		return err
	}

	token, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("generating uuid: %w", err)
	}

	input := &aws_cloudformation_old.DeleteStackInput{
		ClientRequestToken: &token,
		StackName:          stack.StackName,
	}

	if *stack.StackStatus == d.cf.Constants().StackStatusDeleteInProgress {
		_, err := fmt.Fprintf(d.out,
			"stack has status '%s', so let's wait for it to be deleted.\n", *stack.StackStatus)
		if err != nil {
			return err
		}
	} else {
		_, err := d.cf.DeleteStack(input)
		if err != nil {
			return fmt.Errorf("delete Stack: %w", err)
		}
	}

	stackStatus, err := d.waitForDeleteNotInProgress(input)
	if err != nil {
		return fmt.Errorf("waiting for Stack to be deleted: %w", err)
	}

	if stackStatus != d.cf.Constants().StackStatusDeleteComplete {
		return fmt.Errorf("unable to delete stack '%s'. Delete status was '%s'. "+
			"You need to manually fix whatever is blocking this Stack to be deleted. Then run this "+
			"program again", stack, stackStatus)
	}

	return nil
}

func (d *DeleterOld) waitForDeleteNotInProgress(input *aws_cloudformation_old.DeleteStackInput) (string, error) {
	wait := true
	var stack *aws_cloudformation_old.Stack
	var err error

	for i := 0; wait; i++ {
		stack, err = d.getStack(input)
		if err != nil {
			return "", fmt.Errorf("get Stack: %w", err)
		}

		if stack == nil {
			return d.cf.Constants().StackStatusDeleteComplete, nil
		}

		wait = *stack.StackStatus == d.cf.Constants().StackStatusDeleteInProgress

		sleepDuration := d.nextSleep(i)

		_, err := fmt.Fprintf(d.out, "Waiting %d seconds to see if stack deletion is done...\n", sleepDuration)
		if err != nil {
			return "", err
		}

		time.Sleep(time.Second * d.nextSleep(i))
	}

	return *stack.StackStatus, nil
}

func (d *DeleterOld) getStack(input *aws_cloudformation_old.DeleteStackInput) (*aws_cloudformation_old.Stack, error) {
	s := &aws_cloudformation_old.DescribeStacksInput{
		StackName: input.StackName,
	}

	stackResponse, err := d.cf.DescribeStacks(s)
	if err != nil {
		awsError, ok := err.(awserr.RequestFailure)
		if !ok {
			return nil, fmt.Errorf("describe stack: %w", err)
		}

		if strings.Contains(awsError.Message(), "does not exist") {
			return nil, nil
		}

		return nil, fmt.Errorf("describe stack: %w", err)
	}

	if len(stackResponse.Stacks) > 1 {
		return nil, errors.New("internal error, expected 1 stack")
	}

	stack := stackResponse.Stacks[0]

	return stack, nil
}

func (_ *DeleterOld) nextSleep(i int) time.Duration {
	if i == 0 {
		return 3
	}

	if i > 0 && i < 3 {
		return 6
	}

	return 10
}
