package validation

import "testing"

type Te struct {
	s string
}
type IntStruct struct {
	Test *Te
	Age  int `json:"age" between:"20-100"`
}

func TestValidation_IntField(t *testing.T) {
	testCases := []struct {
		scenario string
		is       *IntStruct
		expected bool
	}{
		{
			scenario: "equal to min range",
			is: &IntStruct{
				Age:  20,
				Test: &Te{s: "hello"},
			},
			expected: true,
		},
		{
			scenario: "between min-max range",
			is:       &IntStruct{Age: 55},
			expected: true,
		},
		{
			scenario: "equal to max range",
			is:       &IntStruct{Age: 100},
			expected: true,
		},
		{
			scenario: "smaller than min range",
			is:       &IntStruct{Age: 10},
			expected: false,
		},
		{
			scenario: "greater than max range",
			is:       &IntStruct{Age: 120},
			expected: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.scenario, func(t *testing.T) {
			if err := Validate(tC.is); err != nil && tC.expected {
				t.Errorf("failed, err: %v", err)
			}
		})
	}
}
