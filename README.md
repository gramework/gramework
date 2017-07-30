# gramework [![codecov](https://codecov.io/gh/gramework/gramework/branch/master/graph/badge.svg)](https://codecov.io/gh/gramework/gramework) [![Build Status](https://travis-ci.org/gramework/gramework.svg?branch=master)](https://travis-ci.org/gramework/gramework)

The Good Framework

### Useful links and info

[GoDoc](https://godoc.org/github.com/gramework/gramework)

[Gophers Slack #gramework channel](https://gophers.slack.com)

[Discord Server](https://discord.gg/HkW8DsD)

### What is it?

Fast, highly effective and go-way web framework. You get the simple yet powerful API, we handle optimizations internally.
We glad to see your feature requests and PRs, that are implemented as fast as possible while keeping framework high quality.
SPA-first, so [template engine support is WIP](https://github.com/gramework/gramework/issues/5).

### Project history and "Why?"

Basically, before I've started the project, I need a simple, powerful framework with fair license policy.
First I consulted with lawyers, which license to choose, based on the list of packages that I need to use.
Next, we discussed what to do in order to do everything as correctly as possible.

In our days, `net/http`-based projects are slow and cost-ineffective, so I just write the basic version.

**But.**

Those support HTTP/2, but theoretically we can make it work even with fasthttp.

Those also support websockets, but this is already was done.

**But.** Again.

All our company's solutions are based on fasthttp, so we can use our already stable, optimized solutions.

We can provide stable, faster and more effective functionality with really simple API.

We can support `net/http` handlers with compatibility layer.

We can support multiple handler signature, allow runtime route registration etc.

And even more `We can`.

---

So - **why you may want to use it?**

- Gramework is battle-tested
- Gramework is one of the rare frameworks that can help you serve up to 800k rps even on a 4Gb RAM/i5@2.9GHz/2x1Gbit server
- Gramework make your projects' infrastructure costs more effective by using as less memory as possible
- Gramework helps you serve requests faster, and so it helps you increase conversions ([source 1](https://blog.kissmetrics.com/speed-is-a-killer/), [source 2](https://blog.hubspot.com/marketing/page-load-time-conversion-rates))
- You can build software faster with simple API
- You can achieve agile support and get answers to your questions
- You can just ask a feature and most likely it will be implemented and built in
- You can contact me and donate for high priority feature
- You can be sure that all license questions are OK with gramework
- You can buy a corporate-grade support

### API status

Stable, but not frozen: we adding functions, packages or optional arguments, so you can use new features, but we never break your projects.

**Go >= 1.8 is supported.**

Please, fire an issue or pull request if you want any feature, you find a bug or know how to optimize gramework even more.

Using Gramework with `dep` is highly recommended.

# TOC

- [Benchmarks](#benchmarks)
- [3rd-party license info](#3rd-party-license-info)
- [Basic usage](#basic-usage)
  - [Hello world](#hello-world)
  - [Serving a dir](#serving-a-dir)
  - [Serving prepared bytes](#serving-prepared-bytes)
  - [Using dynamic handlers, part 1](#using-dynamic-handlers-part-1)
  - [Using dynamic handlers, part 2](#using-dynamic-handlers-part-2)
  - [Using dynamic handlers, part 3](#using-dynamic-handlers-part-3)
  - [Using dynamic handlers, part 4](#using-dynamic-handlers-part-4)
  - [Using dynamic handlers, part 5](#using-dynamic-handlers-part-5)

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

### Hello world

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

### Serving a dir

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

### Serving prepared bytes

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
		m := map[string]interface{}{
			"name": "Grame",
			"age": 20,
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
		m := map[string]interface{}{
			"name": "Grame",
			"age": 20,
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

		m := map[string]interface{}{
			"name": "Grame",
			"age": 20,
		}

		if err := ctx.JSON(m); err != nil {
			ctx.Err500()
		}
	})

	app.ListenAndServe()
}
```

### Using dynamic handlers, part 4

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

### Using dynamic handlers, part 5

The example below shows how you can get fasthttp.RequestCtx from gramework.Context and after that it do the same that in part 3:

```go
package main

import (
	"github.com/gramework/gramework"
)

func main() {
	app := gramework.New()

	app.GET("/someJSON", func(ctx *gramework.Context) {
		// same as ctx.WriteString("another data")
		ctx.RequestCtx.WriteString("another data")
	})

	app.ListenAndServe()
}
```
