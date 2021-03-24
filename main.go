package main

import (
	"fmt"
	"github.com/oslokommune/aws-delete-stacks/cmd"
	"os"
)

func main() {
	c := cmd.BuildRootCommand()

	if err := c.Execute(); err != nil {
		fmt.Printf("Error executing program: %s\n", err)
		os.Exit(1)
	}
}
