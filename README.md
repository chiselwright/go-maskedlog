# go-maskedlog

![goversion](https://img.shields.io/github/go-mod/go-version/chiselwright/go-maskedlog) [![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/chiselwright/go-maskedlog) [![Coverage](http://gocover.io/_badge/github.com/chiselwright/go-maskedlog)](http://gocover.io/github.com/chiselwright/go-maskedlog) [![reposize](https://img.shields.io/github/repo-size/chiselwright/go-maskedlog)](https://godoc.org/github.com/chiselwright/go-maskedlog) [![openissues](https://img.shields.io/github/issues/chiselwright/go-maskedlog)](https://github.com/chiselwright/go-maskedlog/issues) [![GitHub pull requests](https://img.shields.io/github/issues-pr/chiselwright/go-maskedlog)](https://github.com/chiselwright/go-maskedlog/pulls) [![Total alerts](https://img.shields.io/lgtm/alerts/g/chiselwright/go-maskedlog.svg?logo=lgtm&logoWidth=18)](https://lgtm.com/projects/g/chiselwright/go-maskedlog/alerts/) [![GitHub language count](https://img.shields.io/github/languages/count/chiselwright/go-maskedlog)](https://github.com/chiselwright/go-maskedlog) ![Maintenance](https://img.shields.io/maintenance/yes/2021) [![GitHub last commit](https://img.shields.io/github/last-commit/chiselwright/go-maskedlog)](https://github.com/chiselwright/go-maskedlog/commits)

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
[![Twitter Follow](https://img.shields.io/twitter/follow/chizcw?style=social)](https://twitter.com/intent/user?screen_name=chizcw)
