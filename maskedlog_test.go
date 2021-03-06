package maskedlog_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/chiselwright/go-maskedlog"
)

func TestLogDebug(t *testing.T) {
	t.Parallel()

	// TODO work out how to test the output properly, or at all
}

func TestGetSingleton(t *testing.T) {
	/* The Plan:
	   - create a new object
	   - ensure there are no know sensitive words
	   - add a sensitive work
	   - create a different new object
	   - both objects should have the same list of sensitive words
	*/

	var log1, log2 maskedlog.MaskLog
	var sens1, sens2 *maskedlog.MaskStrings

	// if the singleton really works as planned, we might have values from other
	// test runs in here
	log1 = maskedlog.GetSingleton()

	sens1 = log1.SensitiveStrings
	if len(*sens1) > 0 {
		t.Errorf("GetSingleton() failed: want length=0, got length=%+v", len(*sens1))
	}

	log1.AddSensitiveValue("First")
	log2 = maskedlog.GetSingleton()

	// get the list of strings, and deeply compare them
	sens1 = log1.SensitiveStrings
	sens2 = log2.SensitiveStrings
	if !reflect.DeepEqual(sens1, sens2) {
		t.Errorf("GetSingleton() failed: want %+v, got %+v", sens1, sens2)
	}

	// adding to log2 should do the same thing, in reverse .. we should still
	// have the same in both
	log2.AddSensitiveValue("Second")
	sens1 = log1.SensitiveStrings
	sens2 = log2.SensitiveStrings
	if !reflect.DeepEqual(sens1, sens2) {
		t.Errorf("GetSingleton() failed: want %+v, got %+v", sens1, sens2)
	}
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

// This is a coupld of DIRECT tests on the SafeString method
// Usually it won't be called directly as is't part of SanitizeInterfaceValues()
// but weird things can happen there so we're extra cautious here
func TestSafeString(t *testing.T) {
	t.Parallel()

	want := "deadxxxf123"
	got := maskedlog.SafeString("deadbeef123")

	if want != got {
		t.Errorf("SafeString() failed; want %q, got %q", want, got)
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
			input:    "deadbeef124",
			expected: "deadxxxf124",
		},
		{
			input:    "deadxxxx-xxxx-xxxx-xxxx-xxxxxxxx5678",
			expected: "deadxxxx-xxxx-xxxx-xxxx-xxxxxxxx5678",
		},
	}

	mlog := maskedlog.GetSingleton()

	for _, test := range tests {
		// start with a clean slate
		mlog.Reset()

		mlog.AddSensitiveValue(test.input)

		// prepare the values
		var v []interface{}
		v = append(v, fmt.Sprintf("TOKEN: %s", test.input))
		v = append(v, "SOMETHING ELSE: aValue")

		// sanitize the value(s)
		mlog.SanitizeInterfaceValues(v)

		// stringify the values (for convenience)
		result := maskedlog.Stringify(v)

		// set out expectations
		want := fmt.Sprintf("TOKEN: %s SOMETHING ELSE: aValue", test.expected)

		// pass / fail?
		if result != want {
			t.Errorf("stringify failed: got: '%v', want '%v' (%d)", result, want, len(test.input))
		}
	}
}
