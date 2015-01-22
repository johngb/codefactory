package codefactory

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNew(t *testing.T) {

	Convey("Create default new CodeFactory", t, func() {

		wantCF := &CodeFactory{
			num:    "0123456789",
			lower:  "abcdefghijklmnopqrstuvwxyz",
			upper:  "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			custom: "",
			prefix: "",
			suffix: "",
			format: "#aaaa",
		}
		cf := New()

		So(cf, ShouldResemble, wantCF)
	})
}

func TestReadable(t *testing.T) {

	Convey("Create readable new CodeFactory", t, func() {

		wantCF := &CodeFactory{
			num:    "0123456789",
			lower:  "abcdefghijkmnopqrstuvwxyz",
			upper:  "ABCDEFGHJKLMNPQRSTUVWXYZ",
			custom: "",
			prefix: "",
			suffix: "",
			format: "#aaaa",
		}
		cf := NewReadable()

		So(cf, ShouldResemble, wantCF)
	})
}

func TestExclude(t *testing.T) {
	var testCases = []struct {
		desc      string
		input     string
		wantNum   string
		wantLower string
		wantUpper string
		wantErr   bool
	}{
		{
			desc:      "basic exclude from each set",
			input:     "1aZ",
			wantNum:   "023456789",
			wantLower: "bcdefghijklmnopqrstuvwxyz",
			wantUpper: "ABCDEFGHIJKLMNOPQRSTUVWXY",
			wantErr:   false,
		},
		{
			desc:      "basic exclude full set",
			input:     "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			wantNum:   "0123456789",
			wantLower: "abcdefghijklmnopqrstuvwxyz",
			wantUpper: "",
			wantErr:   false,
		},
		{
			desc:      "empty exclude",
			input:     "",
			wantNum:   "0123456789",
			wantLower: "abcdefghijklmnopqrstuvwxyz",
			wantUpper: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			wantErr:   false,
		},
		{
			desc:      "invalid exclude",
			input:     "$*",
			wantNum:   "0123456789",
			wantLower: "abcdefghijklmnopqrstuvwxyz",
			wantUpper: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			wantErr:   true,
		},
	}

	for i, tt := range testCases {
		Convey(fmt.Sprintf("Case # %d: %s", i, tt.desc), t, func() {

			cf := New()
			err := cf.Exclude(tt.input)

			So(cf.num, ShouldResemble, tt.wantNum)
			So(cf.upper, ShouldResemble, tt.wantUpper)
			So(cf.lower, ShouldResemble, tt.wantLower)

			if tt.wantErr {
				So(err, ShouldNotBeNil)
			} else {
				So(err, ShouldBeNil)
			}
		})
	}
}

func TestSetCustom(t *testing.T) {
	var testCases = []struct {
		desc       string
		input      string
		input2     string
		wantCustom string
		wantErr    error
	}{
		{
			desc:       "valid custom set",
			input:      "bc346NñŒ",
			wantCustom: "bc346NñŒ",
			wantErr:    nil,
		},
		{
			desc:       "input has whitespace",
			input:      " bc3 46NñŒ",
			wantCustom: "",
			wantErr:    errWhitespace,
		},
		{
			desc:       "has duplicates",
			input:      "bctevb32",
			wantCustom: "",
			wantErr:    errDuplicates,
		},
		{
			desc:       "successive sets",
			input:      "2345",
			input2:     "abcd",
			wantCustom: "abcd",
			wantErr:    nil,
		},
	}

	for i, tt := range testCases {
		Convey(fmt.Sprintf("Case # %d: %s", i, tt.desc), t, func() {

			cf := New()
			err := cf.SetCustom(tt.input)

			if tt.input2 != "" {
				err = cf.SetCustom(tt.input2)
			}

			So(cf.custom, ShouldResemble, tt.wantCustom)
			So(cf.upper, ShouldResemble, defaultUppercase)
			So(cf.lower, ShouldResemble, defaultLowercase)
			So(cf.num, ShouldResemble, defaultNumbers)
			So(err, ShouldEqual, tt.wantErr)
		})
	}
}

func TestSetFormat(t *testing.T) {
	var testCases = []struct {
		desc       string
		input      string
		input2     string
		wantFormat string
		wantErr    error
	}{
		{
			desc:       "valid format set",
			input:      "#xxxx",
			wantFormat: "#xxxx",
			wantErr:    nil,
		},
		{
			desc:       "input has inivalid whitespace",
			input:      "\t\n #xxxx",
			wantFormat: defaultFormat,
			wantErr:    errInvalidFormat,
		},
		{
			desc:       "has numbers",
			input:      "#xxx2a",
			wantFormat: defaultFormat,
			wantErr:    errInvalidFormat,
		},
		{
			desc:       "successive sets",
			input:      "#aaaa",
			input2:     "#xxxa",
			wantFormat: "#xxxa",
			wantErr:    nil,
		},
		{
			desc:       "invalid letters used",
			input:      "#afaa",
			wantFormat: defaultFormat,
			wantErr:    errInvalidFormat,
		},
		{
			desc:       "uppercase letters used",
			input:      "#aAaa",
			wantFormat: defaultFormat,
			wantErr:    errInvalidFormat,
		},
	}

	for i, tt := range testCases {
		Convey(fmt.Sprintf("Case # %d: %s", i, tt.desc), t, func() {

			cf := New()
			err := cf.SetFormat(tt.input)

			if tt.input2 != "" {
				err = cf.SetFormat(tt.input2)
			}

			So(cf.format, ShouldResemble, tt.wantFormat)
			So(cf.upper, ShouldResemble, defaultUppercase)
			So(cf.lower, ShouldResemble, defaultLowercase)
			So(cf.custom, ShouldResemble, "")
			So(cf.num, ShouldResemble, defaultNumbers)
			So(err, ShouldEqual, tt.wantErr)
		})
	}
}

func TestExtendLetters(t *testing.T) {
	var testCases = []struct {
		desc      string
		input     string
		input2    string
		wantLower string
		wantUpper string
		wantErr   error
	}{
		{
			desc:      "valid mix of all valid upper and lowercase extensions",
			input:     "ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖØÙÚÛÜÝÞ" + "µßàáâãäåæçèéêëìíîïðñòóôõöøùúûüýþÿ",
			wantLower: defaultLowercase + "µßàáâãäåæçèéêëìíîïðñòóôõöøùúûüýþÿ",
			wantUpper: defaultUppercase + "ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖØÙÚÛÜÝÞ",
			wantErr:   nil,
		},
		{
			desc:      "input with duplicates",
			input:     "ÀÀ",
			wantLower: defaultLowercase,
			wantUpper: defaultUppercase,
			wantErr:   errDuplicates,
		},
		{
			desc:      "non Latin1 input",
			input:     "ÀÁÂ您好",
			wantLower: defaultLowercase,
			wantUpper: defaultUppercase,
			wantErr:   errNotLatin1,
		},
		{
			desc:      "input with whitespace",
			input:     "À Á",
			wantLower: defaultLowercase,
			wantUpper: defaultUppercase,
			wantErr:   errWhitespace,
		},
		{
			desc:      "uppercase already exists",
			input:     "ÀD",
			wantLower: defaultLowercase,
			wantUpper: defaultUppercase,
			wantErr:   errAlreadyExist,
		},
		{
			desc:      "lowercase already exists",
			input:     "ñc",
			wantLower: defaultLowercase,
			wantUpper: defaultUppercase,
			wantErr:   errAlreadyExist,
		},
		{
			desc:      "non-letter input",
			input:     "ñ9",
			wantLower: defaultLowercase,
			wantUpper: defaultUppercase,
			wantErr:   errNotLetter,
		},
		{
			desc:      "consecutive extends",
			input:     "ñ",
			input2:    "ß",
			wantLower: defaultLowercase + "ñß",
			wantUpper: defaultUppercase,
			wantErr:   nil,
		},
	}

	for i, tt := range testCases {
		Convey(fmt.Sprintf("Case # %d: %s", i, tt.desc), t, func() {

			cf := New()
			err := cf.ExtendLetters(tt.input)

			if tt.input2 != "" {
				err = cf.ExtendLetters(tt.input2)
			}

			So(cf.format, ShouldResemble, defaultFormat)
			So(cf.upper, ShouldResemble, tt.wantUpper)
			So(cf.lower, ShouldResemble, tt.wantLower)
			So(cf.custom, ShouldResemble, "")
			So(cf.num, ShouldResemble, defaultNumbers)
			So(err, ShouldEqual, tt.wantErr)
		})
	}
}

func TestMaxCodes(t *testing.T) {
	var testCases = []struct {
		desc       string
		format     string
		wantNumber int64
	}{
		{
			desc:       "empty format",
			format:     "# $",
			wantNumber: 0,
		},
		{
			desc:       "format with d's",
			format:     "# dd",
			wantNumber: 10 * 10,
		},
		{
			desc:       "format with x's",
			format:     "# dx",
			wantNumber: 10 * (26 + 26 + 10),
		},
		{
			desc:       "format with l's",
			format:     "# ll",
			wantNumber: 26 * 26,
		},
		{
			desc:       "format with u's",
			format:     "# lu",
			wantNumber: 26 * 26,
		},
		{
			desc:       "format with a's",
			format:     "# da",
			wantNumber: 10 * (26 + 26),
		},
		{
			desc:       "really high code count",
			format:     "xxxxxxxxxxx",
			wantNumber: -1,
		},
	}

	for i, tt := range testCases {
		Convey(fmt.Sprintf("Case # %d: %s", i, tt.desc), t, func() {

			cf := New()
			cf.SetFormat(tt.format)

			So(cf.MaxCodes(), ShouldEqual, tt.wantNumber)
		})
	}
}

func TestGenerate(t *testing.T) {

	Convey("testing with '#xxxx'", t, func() {

		cf := New()
		cf.SetFormat("#xxxx")

		res, err := cf.Generate(1)

		So(err, ShouldBeNil)
		So(isIncludedIn((cf.num+cf.upper+cf.lower), rune(res[0][1])), ShouldBeTrue)
		So(isIncludedIn((cf.num+cf.upper+cf.lower), rune(res[0][2])), ShouldBeTrue)
		So(isIncludedIn((cf.num+cf.upper+cf.lower), rune(res[0][3])), ShouldBeTrue)
		So(isIncludedIn((cf.num+cf.upper+cf.lower), rune(res[0][4])), ShouldBeTrue)

	})

	Convey("testing with '# d'", t, func() {

		cf := New()
		cf.SetFormat("# d")

		res, err := cf.Generate(1)

		So(err, ShouldBeNil)
		So(isIncludedIn((cf.num), rune(res[0][2])), ShouldBeTrue)
	})

	Convey("testing with '$ l'", t, func() {

		cf := New()
		cf.SetFormat("$ l")

		res, err := cf.Generate(1)

		So(err, ShouldBeNil)
		So(isIncludedIn((cf.lower), rune(res[0][2])), ShouldBeTrue)
	})

	Convey("testing with '$ u'", t, func() {

		cf := New()
		cf.SetFormat("$ u")

		res, err := cf.Generate(1)

		So(err, ShouldBeNil)
		So(isIncludedIn((cf.upper), rune(res[0][2])), ShouldBeTrue)
	})

	Convey("testing with '$ a'", t, func() {

		cf := New()
		cf.SetFormat("$ a")

		res, err := cf.Generate(1)

		So(err, ShouldBeNil)
		So(isIncludedIn((cf.upper+cf.lower), rune(res[0][2])), ShouldBeTrue)
	})

	Convey("testing with a hopeful code collision", t, func() {

		cf := New()
		cf.SetFormat("$ dd")
		numcodes := 20

		res, err := cf.Generate(numcodes)

		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, numcodes)
	})

	Convey("testing with a too many code collisions", t, func() {

		cf := New()
		cf.SetFormat("$ dd")
		numcodes := 90

		res, err := cf.Generate(numcodes)

		So(err, ShouldEqual, errMaxRetriesExceeded)
		So(len(res), ShouldEqual, 0)
	})

	Convey("testing with a too many codes requested", t, func() {

		cf := New()
		cf.SetFormat("$ d")
		numcodes := 20

		res, err := cf.Generate(numcodes)

		So(err, ShouldEqual, errTooManyCodes)
		So(len(res), ShouldEqual, 0)
	})

}
