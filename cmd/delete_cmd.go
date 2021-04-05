package cmd

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/aws-delete-stacks/core/domain/delete_stacks"
	"github.com/oslokommune/aws-delete-stacks/core/domain/delete_stacks/cloudformation_api"
	"github.com/spf13/cobra"
)

type deleteStacksOpts struct {
	IncludeFilter string
	ExcludeFilter string
	Force         bool
}

// Validate the inputs
func (o *deleteStacksOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.IncludeFilter, validation.Required),
		validation.Field(&o.IncludeFilter, validation.Length(1, 0)),
		validation.Field(&o.ExcludeFilter, validation.Length(1, 0)),
	)
}

func BuildDeleteCommand(cf cloudformation_api.CloudFormation) *cobra.Command {
	opts := &deleteStacksOpts{}

	cmd := &cobra.Command{
		Use:   "aws-delete-listStacksOutput <INCLUDE FILTER>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete AWS cloudformation listStacksOutput",
		Long: "Delete AWS cloudformation listStacksOutput with names containing the string INCLUDE FILTER " +
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
		"Set filter for which listStacksOutput to subtract from included results (filter method: string contains).")
	flags.BoolVarP(&opts.Force, "force", "f", false,
		"Use this flag to actually delete listStacksOutput. Otherwise nothing is deleted.")

	return cmd
}
