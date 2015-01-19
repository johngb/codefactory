package codegen

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExclude(t *testing.T) {
	var testCases = []struct {
	description string
	input       string
	wantNum     string
	wantLower   string
	wantUpper   string
	wantErr     bool
}{
	{"1aZ", "023456789", "bcdefghijklmnopqrstuvwxyz", "ABCDEFGHIJKLMNOPQRSTUVWXYZ", false},
}
	Convey("Test Exclude", t, func() {
		for _, test := range testCases {
			Convey(test.description, func() {

			cf := New()
			cf.Exclude(test.input)
			So(cf.num, ShouldResemble, test.wantNum)
			So(cf.lower, ShouldResemble, test.wantLower)
			So(cf.upper, ShouldResemble, test.wantUpper)

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
		})
	})
}

// func BenchmarkXXX(b *testing.B) {
// 	b.StopTimer()
// 	for _, test := range testCases {
// 		b.StartTimer()

// 		for i := 0; i < b.N; i++ {
// 			XXX(test.input)
// 		}

// 		b.StopTimer()
// 	}
// }
