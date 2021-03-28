package cloudformation_api

import "fmt"

var stackStatusValues []StackStatus

// StackStatus states the status of a Stack
type StackStatus string

//goland:noinspection GoUnusedGlobalVariable
var (
	StackStatusCreateComplete   = newStackStatus("CREATE_COMPLETE")
	StackStatusDeleteInProgress = newStackStatus("DELETE_IN_PROGRESS")
	StackStatusDeleteComplete   = newStackStatus("DELETE_COMPLETE")
	StackStatusDeleteFailed     = newStackStatus("DELETE_FAILED")
)

func newStackStatus(s string) StackStatus {
	ss := StackStatus(s)
	stackStatusValues = append(stackStatusValues, ss)
	return ss
}

func (c *StackStatus) String() string {
	return string(*c)
}

func ParseStackStatus(s string) (StackStatus, error) {
	for _, v := range stackStatusValues {
		if s == v.String() {
			return StackStatus(v), nil
		}
	}

	return "", fmt.Errorf("no StackStatus exists for value '%s'", s)
}
