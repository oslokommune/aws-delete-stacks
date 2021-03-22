package main

import (
	"aws-stack-deleter/cmd"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello")
	c := cmd.BuildRootCommand()

	if err := c.Execute(); err != nil {
		fmt.Printf("Error executing program: %s\n", err)
		os.Exit(1)
	}
}
