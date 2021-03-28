package delete_stacks

import (
	"fmt"
	"github.com/oslokommune/aws-delete-stacks/core/delete_stacks/cloudformation_api"
	"io"
	"strings"
	"time"
)

type Deleter struct {
	out            io.Writer
	cloudformation cloudformation_api.CloudFormation
}

func NewDeleter(cf cloudformation_api.CloudFormation, out io.Writer) *Deleter {
	return &Deleter{
		cloudformation: cf,
		out:            out,
	}
}

func (d *Deleter) DeleteCloudFormationStacks(includeFilter string, excludeFilter string, force bool) error {
	stacks, err := d.listStacks()
	if err != nil {
		return fmt.Errorf("list stacks: %w", err)
	}

	toDelete, err := d.filter(stacks, includeFilter, excludeFilter)
	if err != nil {
		return fmt.Errorf("filter: %w", err)
	}

	err = d.deleteStacks(toDelete, force)
	if err != nil {
		return fmt.Errorf("delete stacks: %w", err)
	}

	return nil
}

func (d *Deleter) listStacks() ([]*cloudformation_api.Stack, error) {
	statusFilter := []cloudformation_api.StackStatus{
		cloudformation_api.StackStatusDeleteInProgress,
		cloudformation_api.StackStatusCreateComplete,
		cloudformation_api.StackStatusDeleteFailed,
	}

	return d.cloudformation.ListStacks(statusFilter)
}

func (d *Deleter) filter(stacks []*cloudformation_api.Stack, includeFilter string, excludeFilter string) ([]*cloudformation_api.Stack, error) {
	filtered := make([]*cloudformation_api.Stack, 0)

	var add *cloudformation_api.Stack
	for _, stack := range stacks {
		add = nil

		if len(includeFilter) > 0 && strings.Contains(stack.Name, includeFilter) {
			add = stack
		}

		if len(excludeFilter) > 0 && strings.Contains(stack.Name, excludeFilter) {
			add = nil
		}

		if add != nil {
			filtered = append(filtered, stack)
		}
	}

	return filtered, nil
}

func (d *Deleter) deleteStacks(toDelete []*cloudformation_api.Stack, force bool) error {
	if force {
		return d.deleteStacksForce(toDelete)
	} else {
		return d.deleteStacksNoForce(toDelete)
	}
}

func (d *Deleter) deleteStacksForce(stacks []*cloudformation_api.Stack) error {
	_, err := fmt.Fprintf(d.out, "- Deleting %d stacks\n", len(stacks))
	if err != nil {
		return err
	}

	for _, stack := range stacks {
		err = d.deleteStack(stack)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Deleter) deleteStack(stack *cloudformation_api.Stack) error {
	_, err := fmt.Fprintf(d.out, "\nDeleting stack: %+v\n", stack)
	if err != nil {
		return err
	}

	if stack.Status == cloudformation_api.StackStatusDeleteInProgress {
		_, err = fmt.Fprintf(d.out,
			"stack has status '%s', so let's wait for it to be deleted.\n", stack.Status)
		if err != nil {
			return err
		}
	} else {
		err = d.cloudformation.DeleteStack(stack)
		if err != nil {
			return fmt.Errorf("delete Stack: %w", err)
		}
	}

	stackStatus, err := d.waitForDeleteNotInProgress(stack)
	if err != nil {
		return fmt.Errorf("waiting for Stack to be deleted: %w", err)
	}

	if stackStatus != cloudformation_api.StackStatusDeleteComplete {
		return fmt.Errorf("unable to delete stack '%s'. Delete status was '%s'. "+
			"You need to manually fix whatever is blocking this Stack to be deleted. Then run this "+
			"program again", stack, stackStatus)
	}

	return nil
}

func (d *Deleter) waitForDeleteNotInProgress(input *cloudformation_api.Stack) (cloudformation_api.StackStatus, error) {
	wait := true
	var stackStatus cloudformation_api.StackStatus
	var err error

	for i := 0; wait; i++ {
		stackStatus, err = d.cloudformation.GetStackStatus(input)
		if err != nil {
			return "", fmt.Errorf("describe stack: %w", err)
		}

		wait = stackStatus == cloudformation_api.StackStatusDeleteInProgress

		sleepDuration := d.nextSleep(i)

		_, err = fmt.Fprintf(d.out, "Waiting %d seconds to see if stack deletion is done...\n", sleepDuration)
		if err != nil {
			return "", err
		}

		time.Sleep(time.Second * d.nextSleep(i))
	}

	return stackStatus, nil
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

func (d *Deleter) deleteStacksNoForce(stacks []*cloudformation_api.Stack) error {
	_, err := fmt.Fprintf(d.out, "- Would delete %d stack(s)\n", len(stacks))
	if err != nil {
		return err
	}

	for _, stack := range stacks {
		_, err = fmt.Fprintln(d.out, stack)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(d.out, "\nNo CloudFormation stacks were deleted, as you didn't specify the --force flag.\n")
	if err != nil {
		return err
	}

	return nil
}
