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
	validFormatChars = "xdlwupa"

	maxRetriesPercent = 10
	maxRetriesBase    = 4
	maxNumCodes       = 100000000
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
)

type CodeFactory struct {
	num    string
	lower  string
	upper  string
	custom string
	format string
	prefix string
	suffix string
}

func New(options ...func(*CodeFactory) error) (*CodeFactory, error) {
	c := &CodeFactory{
		num:    defaultNumbers,
		lower:  defaultLowercase,
		upper:  defaultUppercase,
		custom: defaultCustom,
		prefix: defaultPrefix,
		suffix: defaultSuffix,
		format: defaultFormat,
	}

	// run through all the setup options in the varidac
	for _, option := range options {
		err := option(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func NewReadable() *CodeFactory {
	c, _ := New() // error is only returned if there are inputs to the options varidac in New()
	c.Exclude(defaultUppercase)
	c.Exclude("l")
	return c
}

func (c *CodeFactory) Exclude(s string) error {
	for _, v := range s {
		switch {
		case unicode.IsDigit(v):
			c.num = strings.Replace(c.num, string(v), "", 1)
		case unicode.IsLower(v):
			c.lower = strings.Replace(c.lower, string(v), "", 1)
		case unicode.IsUpper(v): // upper
			c.upper = strings.Replace(c.upper, string(v), "", 1)
		default:
			return errNotLetter
		}
	}
	return nil
}

func (c *CodeFactory) SetCustom(s string) error {
	if hasWhitespace(s) {
		return errWhitespace
	} else if hasDuplicates(s) {
		return errDuplicates
	}
	c.custom = s
	return nil
}

func (c *CodeFactory) SetFormat(s string) error {
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
	c.format = s
	return nil
}

func (c *CodeFactory) ExtendLetters(s string) error {
	if !allLatin1(s) {
		return errNotLatin1
	} else if hasWhitespace(s) {
		return errWhitespace
	} else if hasDuplicates(s) {
		return errDuplicates
	}

	currentUpper := c.upper
	currentLower := c.lower

	for _, v := range s {

		switch {
		// lowercase
		case unicode.IsLower(v):
			if strings.Contains(c.lower, string(v)) {
				c.lower = currentLower
				return errAlreadyExist
			}
			c.lower += string(v)

		case unicode.IsUpper(v):
			if strings.Contains(c.upper, string(v)) {
				c.upper = currentUpper
				return errAlreadyExist
			}
			c.upper += string(v)

		default:
			c.upper = currentUpper
			c.lower = currentLower
			return errNotLetter

		}
	}
	return nil
}

func (c *CodeFactory) MaxCodes() int64 {

	// lengths for optimisation
	lenX := int64(len(c.num + c.upper + c.lower))
	lenD := int64(len(c.num))
	lenL := int64(len(c.lower))
	lenW := int64(len(c.lower + c.num))
	lenU := int64(len(c.upper))
	lenP := int64(len(c.upper + c.num))
	lenA := int64(len(c.upper + c.lower))

	max := int64(1)

	for _, v := range c.format {
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

func (c *CodeFactory) Generate(num int) ([]string, error) {

	if int64(num) > c.MaxCodes() {
		return []string{}, errTooManyCodes
	}

	// strings to build codes from
	x := c.num + c.upper + c.lower
	d := c.num
	l := c.lower
	w := c.lower + c.num
	u := c.upper
	p := c.upper + c.num
	a := c.upper + c.lower

	// lengths for optimisation
	lenX := len(x)
	lenD := len(d)
	lenL := len(l)
	lenW := len(w)
	lenU := len(u)
	lenP := len(p)
	lenA := len(a)

	res := []string{}
	retries := 0
	maxRetries := (num * maxRetriesPercent / 100) + maxRetriesBase

	for i := 1; i <= num; i++ {

		// result string always starts with a prefix
		r := c.prefix

		for _, v := range c.format {
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

				// default:
				// 	panic("John broke the code! Format should not have passed validation")
				// 	return []string{}, errInvalidFormat
			}

		}
		r += c.suffix

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
