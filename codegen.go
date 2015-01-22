package codegen

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
	validFormatChars = "xdlua"
)

var (
	errDuplicates    = errors.New("can't include duplicates")
	errWhitespace    = errors.New("can't include whitespace")
	errAlreadyExist  = errors.New("can't extend with characters that already exist")
	errNotUpper      = errors.New("not uppercase letters")
	errNotLower      = errors.New("not lower letters")
	errInvalidFormat = errors.New("invalid format character")
	errNotLetter     = errors.New("not a letter")
	errNotLatin1     = errors.New("can only extend with Latin1 letters and digits")
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

func NewReadable() *CodeFactory {
	c := New()
	c.Exclude("IlO")
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

func (c *CodeFactory) Extend(s string) error {
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

func (c *CodeFactory) lenNum() int {
	return len(c.num)
}

func (c *CodeFactory) lenUpper() int {
	return len(c.upper)
}

func (c *CodeFactory) lenLower() int {
	return len(c.lower)
}

func (c *CodeFactory) lenLetter() int {
	total := len(c.lower) + len(c.upper)
	return total
}

func (c *CodeFactory) MakeSingle() string {

	// strings to build codes from
	x := c.num + c.upper + c.lower
	d := c.num
	l := c.lower
	u := c.upper
	a := c.upper + c.lower

	// lengths for optimisation
	lenX := len(x)
	lenD := len(d)
	lenL := len(l)
	lenU := len(u)
	lenA := len(a)

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

		// any uppercase letter
		case 'u':
			r += string(u[rand.Intn(lenU)])

		// any letter (upper or lowercase)
		case 'a':
			r += string(a[rand.Intn(lenA)])

		default:
			panic("John broke the code!")
		}

	}

	r += c.suffix

	return r
}

func hasWhitespace(s string) bool {
	for _, v := range s {
		if unicode.IsSpace(rune(v)) {
			return true
		}
	}
	return false
}

func allUpper(s string) bool {
	for _, v := range s {
		if !unicode.IsUpper(v) {
			return false
		}
	}
	return true
}

func allLower(s string) bool {
	for _, v := range s {
		if !unicode.IsLower(v) {
			return false
		}
	}
	return true
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
