package delete_cloudformation_stacks

import (
	"gotest.tools/assert"
	"testing"
)

func TestFilter(t *testing.T) {
	deleter := NewDeleter(nil)

	testCases := []struct {
		name        string
		input       []*deleteStackInput
		stackFilter string
		expected    []*deleteStackInput
	}{
		{
			name: "Should include filtered",
			input: []*deleteStackInput{{
				RoleARN:   "x",
				StackName: "blabla-myenv-dev",
			}},
			stackFilter: "myenv",
			expected: []*deleteStackInput{{
				RoleARN:   "x",
				StackName: "blabla-myenv-dev",
			}},
		},
		{
			name: "Should not include other stack",
			input: []*deleteStackInput{{
				RoleARN:   "x",
				StackName: "blabla-myenve-dev",
			}},
			stackFilter: "-myenv-",
			expected:    make([]*deleteStackInput, 0),
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// When
			output, err := deleter.filter(tc.input, tc.stackFilter)

			// Then
			assert.NilError(t, err)
			assert.DeepEqual(t, tc.expected, output)
		})
	}

}
