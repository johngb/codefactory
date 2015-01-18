package codegen

import (
	"testing"
)

var testCases = []struct {
	input     string
	expected  int
	expectErr bool
}{
	{"1", 1, false},
}

func TestXXX(t *testing.T) {
	for _, test := range testCases {
		actual, err := XXX(test.input)

		if actual != test.expected {
			t.Fatalf("")
		}

		// if we expect an error and there isn't one
		if test.expectErr && err == nil {
			t.Errorf("Name(%v): expected an error, but error is nil", test.input)
		}
		// if we don't expect an error and there is one
		if !test.expectErr && err != nil {
			t.Errorf("Name(%v): expected no error, but error is %s", test.input, err)
		}
	}
}

func BenchmarkXXX(b *testing.B) {
	b.StopTimer()
	for _, test := range testCases {
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			XXX(test.input)
		}

		b.StopTimer()
	}
}
