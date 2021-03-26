package main

import (
	"github.com/oslokommune/aws-delete-stacks/cmd"
	"github.com/oslokommune/aws-delete-stacks/delete_stacks"
	"os"
)

func main() {
	cf := delete_stacks.NewAWSCloudFormation()
	c := cmd.BuildDeleteCommand(cf)

	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
