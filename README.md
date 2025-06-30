# envsubst
[![GoDoc][godoc-img]][godoc-url]
[![License][license-image]][license-url]
[![Github All Releases][releases-image]][releases]

> Environment variables substitution for Go. see docs [below](#docs)

#### Installation:

##### From binaries
Latest stable `envsubst` [prebuilt binaries for 64-bit Linux, or Mac OS X][releases] are available via Github releases.

###### Linux and MacOS
```sh
curl -L https://github.com/allex/envsubst/releases/download/v1.0.4/envsubst-$(uname -s)-$(uname -m) -o envsubst
chmod +x envsubst
sudo mv envsubst /usr/local/bin
```

###### Windows
Download the latest prebuilt binary from [releases page][releases], or if you have curl installed:
```sh
curl -L https://github.com/allex/envsubst/releases/download/v1.0.4/envsubst.exe
```

##### With go
You can install via `go get` (provided you have installed go):
```console
go get github.com/allex/envsubst/cmd/envsubst
```


#### Using via cli
```sh
envsubst < input.tmpl > output.text
echo 'welcome $HOME ${USER:=a8m}' | envsubst
envsubst -help
```

#### Imposing restrictions
There are several command line flags with which you can control the substitution behavior and cause it to stop with an error code when restrictions are not met. This can be handy if you want to avoid creating e.g. configuration files with unset or empty parameters.
Setting a `-fail-fast` flag in conjunction with either no-unset or no-empty or both will result in a faster feedback loop, this can be especially useful when running through a large file or byte array input, otherwise a list of errors is returned.

**Note:** The `-keep-unset` flag automatically disables `-no-unset` and `-no-empty` restrictions when used, as it preserves undefined variables in their original form rather than treating them as errors.

The flags and their restrictions are: 

|__Option__     | __Meaning__    | __Type__ | __Default__  |
| ------------| -------------- | ------------ | ------------ |
|`-i`  | input file  | `string \| stdin` | `stdin`
|`-o`  | output file | `string \| stdout` |  `stdout`
|`-no-digit`  | do not replace variables starting with a digit, e.g. $1 and ${1} | `flag` |  `false` 
|`-no-unset`  | fail if a variable is not set | `flag` |  `false` 
|`-no-empty`  | fail if a variable is set but empty | `flag` | `false`
|`-keep-unset`  | keep undefined variables as their original text instead of substituting them | `flag` | `false`
|`-fail-fast`  | fails at first occurrence of an error, if `-no-empty` or `-no-unset` flags were **not** specified this is ignored | `flag` | `false`

These flags can be combined to form tighter restrictions. 

#### Using `envsubst` programmatically ?
You can take a look on [`_example/main`](https://github.com/allex/envsubst/blob/master/_example/main.go) or see the example below.
```go
package main

import (
	"fmt"
	"github.com/allex/envsubst"
)

func main() {
    input := "welcome $HOME"
    str, err := envsubst.String(input)
    // ...
    buf, err := envsubst.Bytes([]byte(input))
    // ...
    buf, err := envsubst.ReadFile("filename")
}
```
### Docs
> **📖 [Comprehensive API Documentation](API.md)** - Complete developer guide with examples, best practices, and advanced usage patterns

> api docs here: [![GoDoc][godoc-img]][godoc-url]

|__Expression__     | __Meaning__    |
| ----------------- | -------------- |
|`${var}`           | Value of var (same as `$var`)
|`${var-$DEFAULT}`  | If var not set, evaluate expression as $DEFAULT
|`${var:-$DEFAULT}` | If var not set or is empty, evaluate expression as $DEFAULT
|`${var=$DEFAULT}`  | If var not set, evaluate expression as $DEFAULT
|`${var:=$DEFAULT}` | If var not set or is empty, evaluate expression as $DEFAULT
|`${var+$OTHER}`    | If var set, evaluate expression as $OTHER, otherwise as empty string
|`${var:+$OTHER}`   | If var set, evaluate expression as $OTHER, otherwise as empty string
|`$$var`            | Escape expressions. Result will be `$var`. 

<sub>Most of the rows in this table were taken from [here](http://www.tldp.org/LDP/abs/html/refcards.html#AEN22728)</sub>

### See also

* `os.ExpandEnv(s string) string` - only supports `$var` and `${var}` notations

#### License
MIT

[releases]: https://github.com/allex/envsubst/releases
[releases-image]: https://img.shields.io/github/downloads/allex/envsubst/total.svg?style=for-the-badge
[godoc-url]: https://godoc.org/github.com/allex/envsubst
[godoc-img]: https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge
[license-image]: https://img.shields.io/badge/license-MIT-blue.svg?style=for-the-badge
[license-url]: LICENSE
