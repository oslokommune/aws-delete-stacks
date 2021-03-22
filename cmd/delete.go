package cmd

import "github.com/spf13/cobra"
import deleter "aws-stack-deleter/delete_cloudformation_stacks"

func BuildRootCommand() *cobra.Command {
	var stackFilter string
	var force *bool

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete AWS cloudformation stacks",
		Long:  "Delete AWS cloudformation stacks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleter.DeleteCloudFormationStacks(stackFilter, *force)
		},
	}

	f := cmd.Flags()
	// TODO: Change to positional. Consider not doing regex.
	f.StringVarP(&stackFilter, "search", "s", "",
		"A regex pattern to search for in the cloudformation stack names, when deciding which stacks to delete")
	force = f.Bool("force", false, "Use this flag to actually delete stacks. Otherwise nothing is deleted.")

	return cmd
}
