package cloudformation_api

import (
	"fmt"
	"time"
)

type Stack struct {
	Name         string
	Status       StackStatus
	CreationTime *time.Time
}

func (i *Stack) String() string {
	return fmt.Sprintf("%s (%s)", i.Name, i.CreationTime.Format(time.RFC822))
}
