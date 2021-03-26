package main

import (
	"fmt"
	"github.com/oslokommune/aws-delete-stacks/cmd"
	"github.com/oslokommune/aws-delete-stacks/delete_stacks"
	"os"
)

func main() {
	cf := delete_stacks.NewAWSCloudFormation()
	c := cmd.BuildDeleteCommand(cf)

	if err := c.Execute(); err != nil {
		fmt.Printf("Error executing program: %s\n", err)
		os.Exit(1)
	}
}
