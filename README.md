# gramework [![codecov](https://codecov.io/gh/gramework/gramework/branch/master/graph/badge.svg)](https://codecov.io/gh/gramework/gramework) [![Build Status](https://travis-ci.org/gramework/gramework.svg?branch=master)](https://travis-ci.org/gramework/gramework)

The Good Framework

### API status

Stable, but not frozen: we adding functions, packages or optional arguments, so you can use new features, but we never break your projects.
Please, fire an issue or pull request if you want any feature, you find a bug or know how to optimize gramework even more.
Contribution rules will be added ASAP and now we are working on those too.

# TOC

- [Benchmarks](#benchmarks)
- [3rd-party license info](#3rd-party-license-info)
- [Basic usage](#basic-usage)
  - [Serving static data, part 1](#serving-static-data-part-1)
  - [Serving static data, part 2](#serving-static-data-part-2)
  - [Serving static data, part 3](#serving-static-data-part-3)
  - [Serving static data, part 4](#serving-static-data-part-4)
  - [Using dynamic handlers, part 1](#using-dynamic-handlers-part-1)
  - [Using dynamic handlers, part 2](#using-dynamic-handlers-part-2)
  - [Using dynamic handlers, part 3](#using-dynamic-handlers-part-3)
  - [Using dynamic handlers, part 3](#using-dynamic-handlers-part-3-1)
  - [Using dynamic handlers, part 4](#using-dynamic-handlers-part-4)

# Benchmarks

[![benchmark](https://raw.githubusercontent.com/smallnest/go-web-framework-benchmark/master/benchmark.png)](https://github.com/smallnest/go-web-framework-benchmark)

# 3rd-party license info

- Gramework is now powered by [fasthttp](https://github.com/valyala/fasthttp) and custom fasthttprouter, that is embedded now.
  You can find licenses in `/3rd-Party Licenses/fasthttp` and `/3rd-Party Licenses/fasthttprouter`.
- The 3rd autoTLS implementation, placed in `nettls_*.go`, is an integrated version of
  [caddytls](https://github.com/mholt/caddy/tree/d85e90a7b4c06d1698d0b96b695b05d41833fcd3/caddytls), because using it by simple import isn't an option:
  gramework based on `fasthttp`, that is incompatible with `net/http`.
  In [the commit I based on](https://github.com/mholt/caddy/tree/d85e90a7b4c06d1698d0b96b695b05d41833fcd3), caddy is `Apache-2.0` licensed.
  It's license placed in `/3rd-Party Licenses/caddy`. @mholt [allow us](https://github.com/mholt/caddy/issues/1520#issuecomment-286907851) to copy the code in this repo.

# Basic usage

### Serving static data, part 1

The example below will serve "hello, grameworld" and register flag "bind", that allows you to choose another ip/port that gramework should listen:

```go
package main

import (
	"github.com/gramework/gramework"
)

func main() {
	app := gramework.New()

	app.GET("/", "hello, grameworld")

	app.ListenAndServe()
}
```

### Serving static data, part 2

The example below will serve result of expression "15e10" (`150000000000.000000`) and register flag "bind", that allows you to choose another ip/port that gramework should listen:

```go
package main

import (
	"github.com/gramework/gramework"
)

func main() {
	app := gramework.New()

	app.GET("/", 15e10)

	app.ListenAndServe()
}
```

### Serving static data, part 3

The example below will serve static files from ./files and register flag "bind", that allows you to choose another ip/port that gramework should listen:

```go
package main

import (
	"github.com/gramework/gramework"
)

func main() {
	app := gramework.New()

	app.GET("/*any", app.ServeDir("./files"))

	app.ListenAndServe()
}
```

### Serving static data, part 4

The example below will serve bytes and register flag "bind", that allows you to choose another ip/port that gramework should listen:

```go
package main

import (
	"github.com/gramework/gramework"
)

func main() {
	app := gramework.New()

	app.GET("/*any", []byte("some data"))

	app.ListenAndServe()
}
```

### Using dynamic handlers, part 1

The example below will serve JSON and register flag "bind", that allows you to choose another ip/port that gramework should listen:

```go
package main

import (
	"github.com/gramework/gramework"
)

func main() {
	app := gramework.New()

	app.GET("/someJSON", func(ctx *gramework.Context) {
		// NOTE: the map below stands here to show you
		// that gramework supports deep serialization.
		m := map[string]map[string]map[string]map[string]int{
			"abc": {
				"def": {
					"ghk": {
						"wtf": 42,
					},
				},
			},
		}

		if err := ctx.JSON(m); err != nil {
			ctx.Err500()
		}
	})

	app.ListenAndServe()
}
```

### Using dynamic handlers, part 2

The example below will serve JSON with CORS enabled for all routes and register flag "bind", that allows you to choose another ip/port that gramework should listen:

```go
package main

import (
	"github.com/gramework/gramework"
)

func main() {
	app := gramework.New()

	app.Use(app.CORSMiddleware())

	app.GET("/someJSON", func(ctx *gramework.Context) {
		// NOTE: the map below stands here to show you
		// that gramework supports deep serialization.
		m := map[string]map[string]map[string]map[string]int{
			"abc": {
				"def": {
					"ghk": {
						"wtf": 42,
					},
				},
			},
		}

		if err := ctx.JSON(m); err != nil {
			ctx.Err500()
		}
	})

	app.ListenAndServe()
}
```

### Using dynamic handlers, part 3

The example below will serve JSON with CORS enabled in the handler and register flag "bind", that allows you to choose another ip/port that gramework should listen:

```go
package main

import (
	"github.com/gramework/gramework"
)

func main() {
	app := gramework.New()

	app.GET("/someJSON", func(ctx *gramework.Context) {
		ctx.CORS()

		// NOTE: the map below stands here to show you
		// that gramework supports deep serialization.
		m := map[string]map[string]map[string]map[string]int{
			"abc": {
				"def": {
					"ghk": {
						"wtf": 42,
					},
				},
			},
		}

		if err := ctx.JSON(m); err != nil {
			ctx.Err500()
		}
	})

	app.ListenAndServe()
}
```

### Using dynamic handlers, part 3

The example below will serve a string and register flag "bind", that allows you to choose another ip/port that gramework should listen:

```go
package main

import (
	"github.com/gramework/gramework"
)

func main() {
	app := gramework.New()

	app.GET("/someJSON", func(ctx *fasthttp.RequestCtx) {
		ctx.WriteString("another data")
	})

	app.ListenAndServe()
}
```

### Using dynamic handlers, part 4

The example below shows how you can get fasthttp.RequestCtx from gramework.Context and after that it do the same that in part 3:

```go
package main

import (
	"github.com/gramework/gramework"
)

func main() {
	app := gramework.New()

	app.GET("/someJSON", func(ctx *gramework.Context) {
		ctx.RequestCtx.WriteString("another data")
	})

	app.ListenAndServe()
}
```
