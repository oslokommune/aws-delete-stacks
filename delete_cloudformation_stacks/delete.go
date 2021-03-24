package delete_cloudformation_stacks

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func DeleteCloudFormationStacks(stackFilter string, force bool) error {
	mySession := session.Must(session.NewSession())
	cf := cloudformation.New(mySession)

	deleter := NewDeleter(cf)

	err := deleter.DeleteCloudFormationStacks(stackFilter, force)
	if err != nil {
		return fmt.Errorf("delete cloud formation stacks: %w", err)
	}

	return nil
}
