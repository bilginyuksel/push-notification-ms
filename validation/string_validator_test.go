package validation

import (
	"testing"
)

// Another control for the string
// if betweenStr, ok := tag.Lookup("between"); ok {
// 	sep := strings.Split(betweenStr, "-")
// 	minLimit, _ := strconv.Atoi(sep[0])
// 	maxLimit, _ := strconv.Atoi(sep[1])

// 	if len(value) < minLimit || len(value) > maxLimit {
// 		// between limit could not satisfied
// 		return false
// 	}
// }

type TestString struct {
	Email string `json:"email" blank:"false" pattern:"^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$" between:"10-40"`
}

func TestValidation_StringField(t *testing.T) {
	testCases := []struct {
		st       *TestString
		scenario string
		expected bool
	}{
		{
			st:       &TestString{Email: "test@gmail.com"},
			scenario: "example simple gmail",
			expected: false,
		},
		{
			st:       &TestString{Email: "   "},
			scenario: "blank with a lot of spaces",
			expected: false,
		},
		{
			st:       &TestString{Email: "bilgin.yuk21"},
			scenario: "has dot and digits but not mail sign",
			expected: false,
		},
		{
			st:       &TestString{Email: "bilgin.yuksel96@gmail.com"},
			scenario: "my mail address",
			expected: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.scenario, func(t *testing.T) {
			// if validate fails and expected is true then t.error
			if err := Validate(tC.st); err != nil && tC.expected {
				t.Errorf("failed, err: %v\n", err)
			}
		})
	}
}