package cmd

import (
	"fmt"
	"github.com/oslokommune/aws-delete-stacks/cloudformation_api"
	"github.com/oslokommune/aws-delete-stacks/delete_stacks"
	"github.com/spf13/cobra"
)

func BuildDeleteCommand(cf cloudformation_api.CloudFormation) *cobra.Command {
	opts := &deleteStacksOpts{}

	cmd := &cobra.Command{
		Use:   "aws-delete-stacks <INCLUDE FILTER>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete AWS cloudformation stacks",
		Long: "Delete AWS cloudformation stacks with names containing the string INCLUDE FILTER " +
			"(minus stack containing exclude filter).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.IncludeFilter = args[0]

			err := opts.Validate()
			if err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			deleter := delete_stacks.NewDeleter(cf, cmd.OutOrStdout())

			err := deleter.DeleteCloudFormationStacks(opts.IncludeFilter, opts.ExcludeFilter, opts.Force)
			if err != nil {
				fmt.Printf("Program error: %s\n", err)
			}

			return nil
		},
	}

	flags := cmd.Flags()
	cmd.SilenceUsage = true

	flags.StringVarP(&opts.ExcludeFilter, "exclude", "e", "",
		"Set filter for which stacks to subtract from included results (filter method: string contains).")
	flags.BoolVarP(&opts.Force, "force", "f", false,
		"Use this flag to actually delete stacks. Otherwise nothing is deleted.")

	return cmd
}
