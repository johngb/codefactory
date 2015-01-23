# codefactory for Go
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

The output can be extended and controlled by:
- The letters can be extended with any valid Latin1 letters by using the `codefactory.ExtendLetters` method.
- Any Latin1 letters and numbers can be excluded from the output using the `codefactory.Exclude` method.
- Any UTF-8 prefix can be set to appear before each code using the `codefactory.SetPrefix` method as long as it doesn't contain leading whitespace.
- Any UTF-8 suffix can be set to appear after each code using the `codefactory.SetSuffix` method as long as it doesn't contain trailing whitespace.
- Any valid printable (excluding whitespace) Latin1 characters can be included in a custom set using the `codefactory.SetCustom` method.
- The format of the code may be set using the `codefactory.SetFormat` method. The format is interpreted (from the sets controlled with the previous methods)using the following letter codes:
 - `x` = any number, uppercase, or lowercase letter
 - `d` = any number
 - `l` = any lowercase letter
 - `w` = any lowercase letter or number
 - `u` = any uppercase letter
 - `p` = any uppercase letter or number
 - `a` = any uppercase or lowercase letter
 - `c` = any custom character
 - any punctuation, symbol, or whitespace will be printed in the final code, which makes it possible to generate codes such as: `(0)31 36-72-13`

Once the `CodeFactory` has been set up, simply call the `codefactory.Generate` method passing in the number of unique codes required.  An error will be returned if it's not practical to generate the number of codes given the format and sets specified, or if it exceeds the maximum number of codes, which is currently set at 10,000,000.

[See GoDoc](http://godoc.org/github.com/johngb/codefactory) for further documentation.

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

## Performance

Typical performance figures generated on a laptop are shown below with the end of the benchmark name giving the number of codes generated.  I.E. 1E3 = 1,000 while 1E7 = 10,000,000:

#### Without a prefix and a suffix
```
BenchmarkGenerate1E0	  300000	      4390 ns/op
BenchmarkGenerate1E1	  100000	     16973 ns/op
BenchmarkGenerate1E2	   10000	    135598 ns/op
BenchmarkGenerate1E3	    1000	   1307566 ns/op
BenchmarkGenerate1E4	     100	  12419192 ns/op
BenchmarkGenerate1E5	      10	 137495402 ns/op
BenchmarkGenerate1E6	       1	1670830314 ns/op
BenchmarkGenerate1E7	       1	8274087037 ns/op
```

#### With a prefix and a suffix
```
BenchmarkGeneratePS1E0	  100000	     12287 ns/op
BenchmarkGeneratePS1E1	   50000	     34648 ns/op
BenchmarkGeneratePS1E2	    5000	    276009 ns/op
BenchmarkGeneratePS1E3	    1000	   1716471 ns/op
BenchmarkGeneratePS1E4	     100	  14476866 ns/op
BenchmarkGeneratePS1E5	      10	 142052773 ns/op
BenchmarkGeneratePS1E6	       1	1704807308 ns/op
BenchmarkGeneratePS1E7	       1	8581374368 ns/op
```

## Testing and stability

The code has 100% test coverage, but is fairly new, so use it in production with caution.

## Contribution

Contributions welcome! Please fork the repository and open a pull request with your changes.

You are also welcome to open an issue if there's something you feel can be improved or you have a question.

## License

The MIT License (MIT)

Copyright (c) 2015 John Beckett

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
