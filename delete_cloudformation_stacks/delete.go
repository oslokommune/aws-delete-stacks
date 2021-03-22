package delete_cloudformation_stacks

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func DeleteCloudFormationStacks(stackFilter string, force bool) error {
	fmt.Printf("Deleting stacks (%t) with filter %s\n", force, stackFilter)
	fmt.Println("TODO: Implement this")

	return nil
}

func run() error {
	mySession := session.Must(session.NewSession())

	// Create a CloudFormation client from just a session.
	c := cloudformation.New(mySession)

	output, err := c.ListStacks(&cloudformation.ListStacksInput{})
	if err != nil {
		return fmt.Errorf("listing stacks: %w", err)
	}

	fmt.Println(output)

	return nil
}
