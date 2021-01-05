package maskedlog_test

import (
	"fmt"
	"go-maskedlog"
	"testing"
)

func TestLogDebug(t *testing.T) {
	t.Parallel()

	// TODO work out how to test the output properly, or at all
}

func TestStringify(t *testing.T) {
	t.Parallel()

	var want, result string

	var v []interface{}
	v = append(v, "one")
	v = append(v, "two")

	want = "one two"
	result = maskedlog.Stringify(v)
	if result != want {
		t.Errorf("stringify failed: got: '%v', want '%v'", result, want)
	}

	v = append(v, 3)
	want = "one two 3"
	result = maskedlog.Stringify(v)
	if result != want {
		t.Errorf("stringify failed: got: '%v', want '%v'", result, want)
	}

	v = append(v, "four")
	want = "one two 3 four"
	result = maskedlog.Stringify(v)
	if result != want {
		t.Errorf("stringify failed: got: '%v', want '%v'", result, want)
	}

	v = append(v, []int{5, 6})
	want = "one two 3 four [5 6]"
	result = maskedlog.Stringify(v)
	if result != want {
		t.Errorf("stringify failed: got: '%v', want '%v'", result, want)
	}
}

type sanitizeTest struct {
	input    string
	expected string
}

func TestSanitizeInterfaceValues(t *testing.T) {
	t.Parallel()

	var tests = []sanitizeTest{
		{
			input:    "deadbeef-1234-dead-beef-deaffeed5678",
			expected: "deadxxxx-xxxx-xxxx-xxxx-xxxxxxxx5678",
		},
		{
			input:    "deadbeef12",
			expected: "dexxxxxxx2",
		},
		{
			input:    "deadbee-12",
			expected: "dexxxxx-x2",
		},
		{
			input:    "deadbeef123",
			expected: "deadxxxf123",
		},
		{
			input:    "deadxxxx-xxxx-xxxx-xxxx-xxxxxxxx5678",
			expected: "deadxxxx-xxxx-xxxx-xxxx-xxxxxxxx5678",
		},
	}

	// bad way to do things, but this should tide us over until we refactor and
	// improve things further
	// var origToken = vars.GenieToken

	for _, test := range tests {
		// update the GenieToken value
		// vars.GenieToken = test.input

		// prepare the values
		var v []interface{}
		v = append(v, fmt.Sprintf("TOKEN: %s", test.input))
		v = append(v, "SOMETHING ELSE: aValue")

		// sanitize the value(s)
		maskedlog.SanitizeInterfaceValues(v)

		// stringify the values (for convenience)
		result := maskedlog.Stringify(v)

		// set out expectations
		want := fmt.Sprintf("TOKEN: %s SOMETHING ELSE: aValue", test.expected)

		// pass / fail?
		if result != want {
			t.Errorf("stringify failed: got: '%v', want '%v' (%d)", result, want, len(test.input))
		}
	}

	// reset the token
	// vars.GenieToken = origToken
}
