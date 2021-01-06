# go-maskedlog

![goversion](https://img.shields.io/github/go-mod/go-version/chiselwright/go-maskedlog) [![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/chiselwright/go-maskedlog) [![Coverage](http://gocover.io/_badge/github.com/chiselwright/go-maskedlog)](http://gocover.io/github.com/chiselwright/go-maskedlog)

A logger on built on top of [zerolog](https://github.com/rs/zerolog) that can
mask sensitive values in the output.

## Installation

```sh
go get github.com/chiselwright/go-maskedlog
```

## Getting Started

```go
package main

import (
	"fmt"

	"github.com/chiselwright/go-maskedlog"
)

func main() {
	logger := maskedlog.GetSingleton()

	val := "MySekritWurd"
	logger.AddSensitiveValue(val)

	logger.LogWarn(fmt.Sprintf("Failed to authenticate with password: %q", val))
}
```

will result in something similar to:

```txt
❯ go run .
{"level":"warn","time":"2021-01-06T21:25:19Z","message":"Failed to authenticate with password: \"MySexxxxWurd\""}
```
