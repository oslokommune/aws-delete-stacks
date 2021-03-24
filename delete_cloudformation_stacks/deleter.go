package delete_cloudformation_stacks

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/hashicorp/go-uuid"
	"strings"
	"time"
)

type Deleter struct {
	cf *cloudformation.CloudFormation
}

func NewDeleter(cf *cloudformation.CloudFormation) *Deleter {
	return &Deleter{
		cf: cf,
	}
}

func (d *Deleter) DeleteCloudFormationStacks(stackFilter string, force bool) error {
	output, err := d.listStacks()
	if err != nil {
		return err
	}

	var toDelete []*deleteStackInput
	toDelete = d.toDeleteStackInput(output)

	toDelete, err = d.filter(toDelete, stackFilter)
	if err != nil {
		return fmt.Errorf("filter stack input: %w", err)
	}

	err = d.deleteStacks(toDelete, force)
	if err != nil {
		return fmt.Errorf("delete stacks: %w", err)
	}

	return nil
}

func (d *Deleter) listStacks() ([]*cloudformation.StackSummary, error) {
	outputs := make([]*cloudformation.StackSummary, 0)

	statusDeleteInProgres := cloudformation.StackStatusDeleteInProgress
	statusCreateComplete := cloudformation.StackStatusCreateComplete
	statusDeleteFailed := cloudformation.StackStatusDeleteFailed
	statusFilter := []*string{&statusDeleteInProgres, &statusCreateComplete, &statusDeleteFailed}

	i := 0
	crashProtection := 3
	var nextPageToken *string

	for i < crashProtection {
		var output *cloudformation.ListStacksOutput
		var err error

		if i == 0 {
			output, err = d.cf.ListStacks(&cloudformation.ListStacksInput{
				StackStatusFilter: statusFilter,
			})
		} else {
			output, err = d.cf.ListStacks(&cloudformation.ListStacksInput{
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

type deleteStackInput struct {
	StackName   string
	RoleARN     string
	StackStatus string
	Time        time.Time
}

func (i *deleteStackInput) String() string {
	return fmt.Sprintf("%s (%s)", i.StackName, i.Time.Format(time.RFC822))
}

func (d *Deleter) toDeleteStackInput(output []*cloudformation.StackSummary) []*deleteStackInput {
	input := make([]*deleteStackInput, 0)

	for _, summary := range output {
		input = append(input, &deleteStackInput{
			StackName:   *summary.StackName,
			RoleARN:     *summary.StackId,
			StackStatus: *summary.StackStatus,
			Time:        *summary.CreationTime,
		})
	}

	return input
}

func (d *Deleter) filter(stacks []*deleteStackInput, stackFilter string) ([]*deleteStackInput, error) {
	filtered := make([]*deleteStackInput, 0)

	for _, input := range stacks {
		match := strings.Contains(input.StackName, stackFilter)

		if match {
			filtered = append(filtered, input)
		}
	}

	return filtered, nil
}

func (d *Deleter) deleteStacks(stacks []*deleteStackInput, force bool) error {
	if force {
		fmt.Printf("- Deleting %d stacks\n", len(stacks))
	} else {
		fmt.Printf("- Would delete %d stack\n", len(stacks))
	}

	for _, stack := range stacks {
		if !force {
			fmt.Println(stack)
			continue
		}

		err := d.deleteStack(stack)
		if err != nil {
			return err
		}
	}

	fmt.Printf("\nNo CloudFormation stacks were deleted, as you didn't specify the --force flag.")

	return nil
}

func (d *Deleter) deleteStack(stack *deleteStackInput) error {
	fmt.Printf("\nDeleting stack: %+v\n", stack)

	token, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("generating uuid: %w", err)
	}

	input := &cloudformation.DeleteStackInput{
		ClientRequestToken: &token,
		StackName:          &stack.StackName,
	}

	if stack.StackStatus == cloudformation.StackStatusDeleteInProgress {
		fmt.Printf("Stack has status '%s', so let's wait for it to be deleted.\n", stack.StackStatus)
	} else {
		_, err := d.cf.DeleteStack(input)
		if err != nil {
			return fmt.Errorf("delete stack: %w", err)
		}
	}

	stackStatus, err := d.waitForDeleteNotInProgress(input)
	if err != nil {
		return fmt.Errorf("waiting for stack to be deleted: %w", err)
	}

	if stackStatus != cloudformation.StackStatusDeleteComplete {
		return fmt.Errorf("unable to delete stack '%s'. Delete status was '%s'. "+
			"You need to manually fix whatever is blocking this stack to be deleted. Then run this "+
			"program again", stack, stackStatus)
	}

	return nil
}

func (d *Deleter) waitForDeleteNotInProgress(input *cloudformation.DeleteStackInput) (string, error) {
	wait := true
	var stack *cloudformation.Stack
	var err error

	for i := 0; wait; i++ {
		stack, err = d.getStack(input)
		if err != nil {
			return "", fmt.Errorf("get stack: %w", err)
		}

		if stack == nil {
			return cloudformation.StackStatusDeleteComplete, nil
		}

		wait = *stack.StackStatus == cloudformation.StackStatusDeleteInProgress

		sleepDuration := d.nextSleep(i)
		fmt.Printf("Waiting %d seconds to see if stack deletion is done...\n", sleepDuration)

		time.Sleep(time.Second * d.nextSleep(i))
	}

	return *stack.StackStatus, nil
}

func (_ *Deleter) nextSleep(i int) time.Duration {
	if i == 0 {
		return 3
	}

	if i > 0 && i < 3 {
		return 6
	}

	return 10
}

func (d *Deleter) getStack(input *cloudformation.DeleteStackInput) (*cloudformation.Stack, error) {
	s := &cloudformation.DescribeStacksInput{
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
