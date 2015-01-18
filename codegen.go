package codegen

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

const (
	defaultnumbers   = "0123456789"
	defaultlowercase = "abcdefghijklmnopqrstuvwxyz"
	defaultuppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	defaultFormat    = "#aaaa"
)

var (
	errDuplicates   = errors.New("can't include duplicates")
	errWhitespace   = errors.New("can't include whitespace")
	errAlreadyExist = errors.New("can't extend with characters that already exist")
	errNotUpper     = errors.New("not uppercase letters")
	errNotLower     = errors.New("not lower letters")
	errInvalidFormat = errors.New("invalid format character")
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
		custom: "",
		prefix: "",
		suffix: "",
		format: defaultFormat,
	}
}

func (c *CodeFactory) Exclude(s string) error {
	for _, v := range s {
		switch v {
		case v >= '0' && v <= '9': // number
			c.num = strings.Trim(c.num, v)
		case v >= 'a' && v <= 'z': // lower
			c.lower = strings.Trim(c.lower, v)
		case v >= 'A' && v <= 'Z': // upper
			c.upper = strings.Trim(c.upper, v)
		default:
			return fmt.Errorf("the character %q is not a letter or a number", v)
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
    	if !unicode.IsPunct(v) && !unicode.IsSymbol(v) && != ' ' {
    		return errInvalidFormat
    	}
    }
	c.format = s
	return nil
}

func (c *CodeFactory) ExtendUpper(s string) error {
	if hasWhitespace(s) {
		return errWhitespace
	} else if hasDuplicates(s) {
		return errDuplicates
	} else if hasDuplicates(s + c.upper) {
		return errAlreadyExist
	} else if !allUpper(s) {
		return errNotUpper
	}
	c.upper += s
	return nil
}

func (c *CodeFactory) ExtendLower(s string) error {
	if hasWhitespace(s) {
		return errCustomHasWhitespace
	} else if hasDuplicates(s) {
		return errCustomHasDuplicates
	} else if hasDuplicates(s + c.lower) {
		return errDuplicateOfUpper
	} else if !allLower(s) {
		return errNotLower
	}
	c.upper += s
	return nil
}

func hasWhiteSpace(s string) bool {
	for _, c := range s {
		if unicode.IsSpace(rune(c)) {
			return true
		}
	}
	return false
}

func allUpper(s string) bool {
	for _, c := range s {
		if !unicode.IsUpper(c) {
			return false
		}
	}
	return true
}

func allLower(s string) bool {
	for _, c := range s {
		if !unicode.IsLower(c) {
			return false
		}
	}
	return true
}

// func hasNonASCII(s string) bool {
// 	for _, c := range s {
// 		if rune(c) > unicode.MaxASCII {
// 			return true
// 		}
// 	}
// 	return false
// }

func hasDuplicates(s string) bool {
	for _, c := range s {
		if strings.Count(s, c) != 1 {
			return true
		}
	}
	return false
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

// func unique(s string) bool {
// 	for _, c := range s {
// 		count := strings.Count(s, c); if count != 1 {
// 			return false
// 		}
// 	}
// 	return true
// }

// func isLower(s string) bool {
// 		for _, c := range s {
// 		if s < 'a' || s > 'z' {
// 			return false
// 		}
// 	}
// 	return true
// }

// func isUpper(s string) bool {
// 		for _, c := range s {
// 		if s < 'A' || s > 'Z' {
// 			return false
// 		}
// 	}
// 	return true
// }

// func isNumber(s string) bool {
// 		for _, c := range s {
// 		if s < '0' || s > '9' {
// 			return false
// 		}
// 	}
// 	return true
// }
