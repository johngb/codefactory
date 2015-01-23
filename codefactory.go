package codefactory

import (
	"errors"
	"math/rand"
	"strings"
	"unicode"
)

const (
	defaultNumbers   = "0123456789"
	defaultLowercase = "abcdefghijklmnopqrstuvwxyz"
	defaultUppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	defaultFormat    = "#aaaa"
	defaultCustom    = ""
	defaultPrefix    = ""
	defaultSuffix    = ""
	validFormatChars = "xdlwupac"

	maxRetriesPercent = 10
	maxRetriesBase    = 4
	maxNumCodes       = 1E6
)

var (
	errDuplicates         = errors.New("can't include duplicates")
	errWhitespace         = errors.New("can't include whitespace")
	errAlreadyExist       = errors.New("can't extend with characters that already exist")
	errNotUpper           = errors.New("not uppercase letters")
	errNotLower           = errors.New("not lower letters")
	errInvalidFormat      = errors.New("invalid format character")
	errNotLetter          = errors.New("not a letter")
	errNotLatin1          = errors.New("can only extend with Latin1 letters and digits")
	errMaxRetriesExceeded = errors.New("too many duplicate codes generated. Consider using a longer code")
	errTooManyCodes       = errors.New("too many codes to generate with given settings")
	errLeadingWhitespace  = errors.New("a prefix may not have leading whitespace")
	errTrailingWhitespace = errors.New("a suffix may not have trailing whitespace")
	errNoCharacters       = errors.New("no characters can be generated with an empty set")
)

// CodeFactory is the type to hold all the config settings to generate codes.
type CodeFactory struct {
	num    string
	lower  string
	upper  string
	custom string
	format string
	prefix string
	suffix string
}

// New generates a new default CodeFactory.
func New() *CodeFactory {
	return &CodeFactory{
		num:    defaultNumbers,
		lower:  defaultLowercase,
		upper:  defaultUppercase,
		custom: defaultCustom,
		prefix: defaultPrefix,
		suffix: defaultSuffix,
		format: defaultFormat,
	}
}

// NewReadable generates a Codefactory with more easily readable default
// codes.  It excludes all uppercase characters, and removes lowercase letters
// that are easily confused with numbers.
func NewReadable() *CodeFactory {
	cf := New() // error is only returned if there are inputs to the options varidac in New()
	cf.Exclude(defaultUppercase)
	cf.Exclude("l")
	return cf
}

// Exclude excludes all characters in the input string from either the
// uppercase, lowercase, or numbers sets. It does not affect the prefix,
// suffix, or custom set.
func (cf *CodeFactory) Exclude(s string) error {
	for _, v := range s {
		switch {
		case unicode.IsDigit(v):
			cf.num = strings.Replace(cf.num, string(v), "", 1)
		case unicode.IsLower(v):
			cf.lower = strings.Replace(cf.lower, string(v), "", 1)
		case unicode.IsUpper(v): // upper
			cf.upper = strings.Replace(cf.upper, string(v), "", 1)
		default:
			return errNotLetter
		}
	}
	return nil
}

// SetCustom sets the custom set of characters.
func (cf *CodeFactory) SetCustom(s string) error {
	if hasWhitespace(s) {
		return errWhitespace
	} else if hasDuplicates(s) {
		return errDuplicates
	}
	cf.custom = s
	return nil
}

// SetFormat sets the format of the codes to be generated.
//
//  x = any number, uppercase, or lowercase letter
//  d = any number
//  l = any lowercase letter
//  w = any lowercase letter or number
//  u = any uppercase letter
//  p = any uppercase letter or number
//  a = any uppercase or lowercase letter
//  c = any character in the custom set
//
// Other than the characters given, the format string may include symbols,
// spaces, and punctuation.
//
// If letters are needed, they should be set in the prefix or suffix.
func (cf *CodeFactory) SetFormat(s string) error {
	for _, v := range s {
		// if not punctuation, symbol, or space
		if !unicode.IsPunct(v) && !unicode.IsSymbol(v) && v != ' ' && !unicode.IsLetter(v) {
			return errInvalidFormat
		}
		if unicode.IsLetter(v) {
			if !isIncludedIn(validFormatChars, v) {
				return errInvalidFormat
			}
		}
	}
	cf.format = s
	return nil
}

// ExtendLetters allows the uppercase and lowercase letter to be extended with
// letters that are part of the Latin1 set of letters.  This allows the
// addition of common letters from Latin script based character sets such as
// German, Spanish, Hungarian, and Norwegian.
func (cf *CodeFactory) ExtendLetters(s string) error {
	if !allLatin1(s) {
		return errNotLatin1
	} else if hasWhitespace(s) {
		return errWhitespace
	} else if hasDuplicates(s) {
		return errDuplicates
	}

	currentUpper := cf.upper
	currentLower := cf.lower

	for _, v := range s {

		switch {
		// lowercase
		case unicode.IsLower(v):
			if strings.Contains(cf.lower, string(v)) {
				cf.lower = currentLower
				return errAlreadyExist
			}
			cf.lower += string(v)

		case unicode.IsUpper(v):
			if strings.Contains(cf.upper, string(v)) {
				cf.upper = currentUpper
				return errAlreadyExist
			}
			cf.upper += string(v)

		default:
			cf.upper = currentUpper
			cf.lower = currentLower
			return errNotLetter

		}
	}
	return nil
}

// SetPrefix allows any UTF-8 prefix to be set before the code.  However, it
// doesn't allow any leading whitespace.
//
// It's possible to use this go generate codes such as:
//  红 88
//  ரெட் 88
//  สีแดง 88
func (cf *CodeFactory) SetPrefix(s string) error {

	// if the prefix is being cleared
	if len(s) == 0 {
		cf.prefix = ""
		return nil
	}

	// else check that there is no leading whitespace
	if unicode.IsSpace(rune(s[0])) {
		return errLeadingWhitespace
	}
	cf.prefix = s
	return nil
}

// SetSuffix allows any UTF-8 suffix to be set after the prefix and code.
// However, it doesn't allow any trailing whitespace.
//
// It's possible to use this go generate codes such as:
//  abc 红
//  abc ரெட்
//  abc สีแดง
func (cf *CodeFactory) SetSuffix(s string) error {

	// if the suffix is being cleared
	if len(s) == 0 {
		cf.suffix = ""
		return nil
	}

	// else check that there is no trailing whitespace
	if unicode.IsSpace(rune(s[len(s)-1])) {
		return errTrailingWhitespace
	}
	cf.suffix = s
	return nil
}

// MaxCodes returns the maximum number of codes that could theoretically be
// generated with the current CodeFactory settings.
//
// In general it will not be possible to generate this full set using
// codefactory as this would cause too many collisions, and require different
// logic to complete.
func (cf *CodeFactory) MaxCodes() int64 {

	// lengths for optimisation
	lenX := int64(len(cf.num + cf.upper + cf.lower))
	lenD := int64(len(cf.num))
	lenL := int64(len(cf.lower))
	lenW := int64(len(cf.lower + cf.num))
	lenU := int64(len(cf.upper))
	lenP := int64(len(cf.upper + cf.num))
	lenA := int64(len(cf.upper + cf.lower))
	lenC := int64(len(cf.custom))

	max := int64(1)

	for _, v := range cf.format {
		if unicode.IsLower(v) {
			switch v {
			case 'x': // any
				max *= lenX
			case 'd': // digits
				max *= lenD
			case 'l': // lowercase
				max *= lenL
			case 'w': // lowercase + number
				max *= lenW
			case 'u': // uppercase
				max *= lenU
			case 'p': // uppercase + number
				max *= lenP
			case 'a': // lowercase + uppercase
				max *= lenA
			case 'c': // custom
				max *= lenC
				// default:
				// 	panic("Invalid format was passed.  Format code possibly broken.")
			}
		}
		if max > maxNumCodes*10 { // ten times the max number of codes to be sure that we can easily generate the given number
			return -1
		}
	}
	if max == 1 {
		return 0
	}
	return max

}

// Generate generates 'num' codes using the settings given in 'cf', and returns
// the codes as an unordered slice of strings.  It will return an error if the
// number of codes is too hight for the given format and character sets in
// 'cf', or if 'num' is greater than the maximum allowed, which is currently
// set at 1,000,000 codes.
func (cf *CodeFactory) Generate(num int) ([]string, error) {

	maxCodes := cf.MaxCodes()
	if maxCodes == 0 {
		return []string{}, errNoCharacters
	} else if int64(num) > cf.MaxCodes() {
		return []string{}, errTooManyCodes
	}

	// strings to build codes from
	x := cf.num + cf.upper + cf.lower
	d := cf.num
	l := cf.lower
	w := cf.lower + cf.num
	u := cf.upper
	p := cf.upper + cf.num
	a := cf.upper + cf.lower
	c := cf.custom

	// lengths for optimisation
	lenX := len(x)
	lenD := len(d)
	lenL := len(l)
	lenW := len(w)
	lenU := len(u)
	lenP := len(p)
	lenA := len(a)
	lenC := len(c)

	res := []string{}
	retries := 0
	maxRetries := (num * maxRetriesPercent / 100) + maxRetriesBase

	for i := 1; i <= num; i++ {

		// result string always starts with a prefix
		r := cf.prefix

		for _, v := range cf.format {
			// formatting symbol
			if !unicode.IsLetter(v) {
				r += string(v)
				continue
			}
			// is code character
			switch v {

			// any character in sets
			case 'x':
				r += string(x[rand.Intn(lenX)])

			// any number digit
			case 'd':
				r += string(d[rand.Intn(lenD)])

			// any lowercase letter
			case 'l':
				r += string(l[rand.Intn(lenL)])

			// any lowercase letter or number
			case 'w':
				r += string(w[rand.Intn(lenW)])

			// any uppercase letter
			case 'u':
				r += string(u[rand.Intn(lenU)])

			// any uppercase letter or number
			case 'p':
				r += string(p[rand.Intn(lenP)])

			// any letter (upper or lowercase)
			case 'a':
				r += string(a[rand.Intn(lenA)])

				// any custom character
			case 'c':
				r += string(c[rand.Intn(lenC)])

				// default:
				// 	panic("John broke the code! Format should not have passed validation")
				// 	return []string{}, errInvalidFormat
			}

		}
		r += cf.suffix

		// check if r is in res
		if exist := isInSlice(res, r); exist {
			i-- // generate a new code
			retries++
			if retries > maxRetries {
				return []string{}, errMaxRetriesExceeded
			}

			continue
		}

		res = append(res, r)
	}
	return res, nil
}

func isInSlice(s []string, r string) bool {
	for _, v := range s {
		if r == v {
			return true
		}
	}
	return false
}

func hasWhitespace(s string) bool {
	for _, v := range s {
		if unicode.IsSpace(rune(v)) {
			return true
		}
	}
	return false
}

func hasDuplicates(s string) bool {
	for _, v := range s {
		if strings.Count(s, string(v)) != 1 {
			return true
		}
	}
	return false
}

func allLatin1(s string) bool {
	for _, v := range s {
		if v > unicode.MaxLatin1 {
			return false
		}
	}
	return true
}

func isIncludedIn(s string, v rune) bool {
	for _, n := range s {
		if v == n {
			return true
		}
	}
	return false
}
