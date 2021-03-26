package cmd

import validation "github.com/go-ozzo/ozzo-validation/v4"

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
