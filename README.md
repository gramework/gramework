# gramework  [![codecov](https://codecov.io/gh/gramework/gramework/branch/master/graph/badge.svg)](https://codecov.io/gh/gramework/gramework) [![Build Status](https://travis-ci.org/gramework/gramework.svg?branch=master)](https://travis-ci.org/gramework/gramework) [![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/1203/badge)](https://bestpractices.coreinfrastructure.org/projects/1203) [![Backers on Open Collective](https://opencollective.com/gramework/backers/badge.svg)](#backers) [![Sponsors on Open Collective](https://opencollective.com/gramework/sponsors/badge.svg)](#sponsors)

The Good Framework

[![Gramework Stats Screenshot](https://10357-5.s.cdn13.com/docs/gramework_stats_screenshot.png)](https://grafana.com/dashboards/3422)

_Gramework long-term testing stand metrics screenshot made with [Gramework Stats Dashboard](https://grafana.com/dashboards/3422) and [metrics middleware](https://github.com/gramework/gramework/tree/dev/metrics)_

### Useful links and info

If you find it, you can submit vulnerability via k@gramework.win.

| Name  | Link/Badge  	|
|---	|---		|
| Docs  | [GoDoc](https://godoc.org/github.com/gramework/gramework) |
| Our Jira | [Jira](https://gramework.atlassian.net) |
| Support us with a donation or become a sponsor | [OpenCollective](https://opencollective.com/gramework) |
| We have #gramework channel in the Gophers Slack | https://gophers.slack.com |
| Our Discord Server | https://discord.gg/HkW8DsD |
| Master branch coverage | [![codecov](https://codecov.io/gh/gramework/gramework/branch/master/graph/badge.svg)](https://codecov.io/gh/gramework/gramework) |
| Master branch status | [![Build Status](https://travis-ci.org/gramework/gramework.svg?branch=master)](https://travis-ci.org/gramework/gramework) |
| Dev branch coverage | [![codecov](https://codecov.io/gh/gramework/gramework/branch/dev/graph/badge.svg)](https://codecov.io/gh/gramework/gramework) |
| Dev branch status | [![Build Status](https://travis-ci.org/gramework/gramework.svg?branch=dev)](https://travis-ci.org/gramework/gramework) |
| CII Best Practices | [![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/1203/badge)](https://bestpractices.coreinfrastructure.org/projects/1203) |
| Gramework Stats Dashboard for Grafana | https://grafana.com/dashboards/3422 |

### What is it?

Fast, highly effective and go-way web framework. You get the simple yet powerful API, we handle optimizations internally.
We glad to see your feature requests and PRs, that are implemented as fast as possible while keeping framework high quality.
SPA-first, so [template engine support is WIP](https://github.com/gramework/gramework/issues/5).

### Gramework is trusted by such projects as:

[![Confideal banner](https://10357-5.s.cdn13.com/docs/confideal_banner.jpg)](https://confideal.io)

Confideal is a service for making fast and safe international deals through smart contracts on Ethereum blockchain.

> With Gramework, we have made a number of improvements:
> - reduced boilerplate code;
> - expedited the development process without cutting neither scope nor performance requirements;
> - reduced the code that needs to be maintained;
> - saved hundreds of hours by using by using functionality that comes as a part of the framework;
> - optimized costs  of service maintenance;
> - took advantage of services' and the implementation code's being scalable.

-----


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

## Contributors

This project exists thanks to all the people who contribute. [[Contribute](CONTRIBUTING.md)].
<a href="graphs/contributors"><img src="https://opencollective.com/gramework/contributors.svg?width=890&button=false" /></a>


## Backers

Thank you to all our backers! 🙏 [[Become a backer](https://opencollective.com/gramework#backer)]

<a href="https://opencollective.com/gramework#backers" target="_blank"><img src="https://opencollective.com/gramework/backers.svg?width=890"></a>


## Sponsors

Support this project by becoming a sponsor. Your logo will show up here with a link to your website. [[Become a sponsor](https://opencollective.com/gramework)]

<a href="https://opencollective.com/gramework/sponsor/0/website" target="_blank"><img src="https://opencollective.com/gramework/sponsor/0/avatar.svg"></a>
<a href="https://opencollective.com/gramework/sponsor/1/website" target="_blank"><img src="https://opencollective.com/gramework/sponsor/1/avatar.svg"></a>
<a href="https://opencollective.com/gramework/sponsor/2/website" target="_blank"><img src="https://opencollective.com/gramework/sponsor/2/avatar.svg"></a>
<a href="https://opencollective.com/gramework/sponsor/3/website" target="_blank"><img src="https://opencollective.com/gramework/sponsor/3/avatar.svg"></a>
<a href="https://opencollective.com/gramework/sponsor/4/website" target="_blank"><img src="https://opencollective.com/gramework/sponsor/4/avatar.svg"></a>
<a href="https://opencollective.com/gramework/sponsor/5/website" target="_blank"><img src="https://opencollective.com/gramework/sponsor/5/avatar.svg"></a>
<a href="https://opencollective.com/gramework/sponsor/6/website" target="_blank"><img src="https://opencollective.com/gramework/sponsor/6/avatar.svg"></a>
<a href="https://opencollective.com/gramework/sponsor/7/website" target="_blank"><img src="https://opencollective.com/gramework/sponsor/7/avatar.svg"></a>
<a href="https://opencollective.com/gramework/sponsor/8/website" target="_blank"><img src="https://opencollective.com/gramework/sponsor/8/avatar.svg"></a>
<a href="https://opencollective.com/gramework/sponsor/9/website" target="_blank"><img src="https://opencollective.com/gramework/sponsor/9/avatar.svg"></a>



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

The example below will serve "hello, grameworld". Gramework will register flag "bind" for you, that allows you to choose another ip/port that gramework should listen:

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

If you don't want support `bind` flag, pass the optional address argument to `ListenAndServe`.

**NOTE**: all examples below will register `bind` flag.

### JSON world ;) Part 1

From version: 1.1.0-rc1

The example below will serve `{"hello":"grameworld"}` from the map. Gramework will register flag "bind" for you, that allows you to choose another ip/port that gramework should listen:

```go
package main

import (
	"github.com/gramework/gramework"
)

func main() {
	app := gramework.New()

	app.GET("/", func() map[string]interface{} {
		return map[string]interface{}{
			"hello": "gramework",
		}
	})

	app.ListenAndServe()
}
```

### JSON world. Part 2

From version: 1.1.0-rc1

The example below will serve `{"hello":"grameworld"}` from the struct. Gramework will register flag "bind" for you, that allows you to choose another ip/port that gramework should listen:

```go
package main

import (
	"github.com/gramework/gramework"
)

type SomeResponse struct {
	hello string
}

func main() {
	app := gramework.New()

	app.GET("/", func() interface{} {
		return SomeResponse{
			hello: "gramework",
		}
	})

	app.ListenAndServe()
}
```

### Serving a dir

The example below will serve static files from ./files:

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

The example below will serve byte slice:

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

### Using dynamic handlers, part 1. Simple JSON response.

The example below will serve JSON:

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

### Using dynamic handlers, part 2. Simple JSON response with service-wide CORS enabled.

The example below will serve JSON with CORS enabled for all routes:

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

### Using dynamic handlers, part 3. Simple JSON response with handler-wide CORS enabled.

The example below will serve JSON with CORS enabled in the handler:

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

### Using dynamic handlers, part 4. Simple FastHTTP-compatible handlers.

The example below will serve a string:

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

### Using dynamic handlers, part 5. Access to fasthttp.RequestCtx from gramework.Context

The example below shows how you can get fasthttp.RequestCtx from gramework.Context and use it:

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
