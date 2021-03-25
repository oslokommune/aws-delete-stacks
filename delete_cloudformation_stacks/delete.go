package delete_cloudformation_stacks

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func DeleteCloudFormationStacks(stackFilter string, force bool) error {
	t := true
	s := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			CredentialsChainVerboseErrors: &t,
		},
		SharedConfigState: session.SharedConfigEnable,
	}))
	cf := cloudformation.New(s)

	deleter := NewDeleter(cf)

	err := deleter.DeleteCloudFormationStacks(stackFilter, force)
	if err != nil {
		return fmt.Errorf("delete cloud formation stacks: %w", err)
	}

	return nil
}
