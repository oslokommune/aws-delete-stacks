package main

import (
	"github.com/oslokommune/aws-delete-stacks/cloudformation_api"
	"github.com/oslokommune/aws-delete-stacks/cmd"
	"os"
)

func main() {
	cf := cloudformation_api.NewAWSCloudFormation()
	c := cmd.BuildDeleteCommand(cf)

	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
