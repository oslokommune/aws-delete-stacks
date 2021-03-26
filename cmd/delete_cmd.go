package cmd

import (
	"fmt"
	"github.com/oslokommune/aws-delete-stacks/delete_stacks"
	"github.com/spf13/cobra"
)

func BuildDeleteCommand(cf delete_stacks.CloudFormation) *cobra.Command {
	opts := &deleteStacksOpts{}

	cmd := &cobra.Command{
		Use:   "delete <FILTER>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete AWS cloudformation stacks",
		Long:  "Delete AWS cloudformation stacks with names that contains the given FILTER.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.IncludeFilter = args[0]

			err := opts.Validate()
			if err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := delete_stacks.DeleteCloudFormationStacks(cf, cmd.OutOrStdout(), opts.IncludeFilter, opts.Force)
			if err != nil {
				fmt.Printf("Program error: %s\n", err)
			}

			return nil
		},
	}

	flags := cmd.Flags()
	cmd.SilenceErrors = true
	flags.BoolVarP(&opts.Force, "force", "f", false, "Use this flag to actually delete stacks. Otherwise nothing is deleted.")

	return cmd
}
