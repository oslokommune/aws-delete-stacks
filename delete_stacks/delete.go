package delete_stacks

import (
	"fmt"
	"io"
)

func DeleteCloudFormationStacks(cf CloudFormation, out io.Writer, stackFilter string, force bool) error {
	deleter := newDeleter(cf, out)

	err := deleter.DeleteCloudFormationStacks(stackFilter, force)
	if err != nil {
		return fmt.Errorf("delete cloud formation stacks: %w", err)
	}

	return nil
}
