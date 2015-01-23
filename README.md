# codefactory
codefactory is a Go package designed to efficiently generate unique codes while giving control as to the exact format that the codes should take.

## Installation

```bash
go get github.com/johngb/codefactory
```

## Usage
The default CodeFactory generated with `codefactory.New()` contains as its valid sets:
 - all uppercase ASCII letters (A-Z)
 - all lowercase ASCII letters (a-z)
 - all ASCII numbers (0-9)

The letters can be extended with any valid Latin1 letters by using the `codefactory.ExtendLetters` method.

Any Latin1 letters and numbers can be excluded from the output using the `codefactory.Exclude` method.

Any UTF-8 prefix can be set to appear before each code using the `codefactory.SetPrefix` method as long as it doesn't contain leading whitespace.

Any UTF-8 suffix can be set to appear after each code using the `codefactory.SetSuffix` method as long as it doesn't contain trailing whitespace.

Any valid printable (excluding whitespace) Latin1 characters can be included in a custom set using the `codefactory.SetCustom` method.

The format of the code may be set using the `codefactory.SetFormat` method. The format is interpreted (from the sets controlled with the previous methods)using the following letter codes:
 - x = any number, uppercase, or lowercase letter
 - d = any number
 - l = any lowercase letter
 - w = any lowercase letter or number
 - u = any uppercase letter
 - p = any uppercase letter or number
 - a = any uppercase or lowercase letter
 - c = any custom character
 - any punctuation, symbol, or whitespace will be printed in the final code, which makes it possible to generate codes such as: `(0)31 36-72-13`

[See GoDoc](http://godoc.org/github.com/johngb/codefactory) for further documentation and examples.

## Example

Let's say that I want to generate 1,000 codes that follow this code in structure: 	// Let's say that I desire a number of codes that follow this example: `代码 (23 aDn) 是重要的`

```Go
package main

import "github.com/johngb/codefactory"

func main() {

	// create a CodeFactory
	cf := codefactory.New()

	// exclude some characters from the default set
	cf.Exclude("iloO01")

	// add a custom prefix
	cf.SetPrefix("代码 (")

	// add a custom suffix
	cf.SetSuffix(") 是重要的")

	// set a format using the formatting characters
	// d = any number in the number set
	// a = any letter in either the uppercase or lowercase letter sets
	cf.SetFormat("#dd aaa")

	output, err := cf.Generate(1000)
	if err != nil {
		// handle the error
	}

	// do something with output
	// ...
}
```

output:

```
代码 (#83 WWV) 是重要的
代码 (#45 xRq) 是重要的
代码 (#73 Ffb) 是重要的
...
```


