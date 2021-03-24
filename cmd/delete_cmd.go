package cmd

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	deleter "github.com/oslokommune/aws-delete-stacks/delete_cloudformation_stacks"
	"github.com/spf13/cobra"
)

type deleteStacksOpts struct {
	StackFilter string
	Force       bool
}

// Validate the inputs
func (o *deleteStacksOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.StackFilter, validation.Required),
		validation.Field(&o.StackFilter, validation.Length(1, 0)),
	)
}

func BuildRootCommand() *cobra.Command {
	opts := &deleteStacksOpts{}

	cmd := &cobra.Command{
		Use:   "delete <FILTER>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete AWS cloudformation stacks",
		Long:  "Delete AWS cloudformation stacks with names that contains the given FILTER.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.StackFilter = args[0]

			err := opts.Validate()
			if err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := deleter.DeleteCloudFormationStacks(opts.StackFilter, opts.Force)
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
