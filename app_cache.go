// +build cache

package gramework

import (
	"errors"
	"time"

	"github.com/VictoriaMetrics/fastcache"
)

// CacheOptions is a handler cache configuration structure.
type CacheOptions struct {
	// TTL is the time that cached response is valid
	TTL time.Duration
	// Cacheable function returns if current request is cacheable.
	// By deafult, any request with Authentication header or any Cookies will not be cached for security reasons.
	// If you want to cache responses for authorized users, please replace both Cacheable and CacheKey functions
	// to make sure that CacheKey includes something like session id.
	Cacheable func(ctx *Context) bool
	// CacheKey function returns the cache key for current request
	CacheKey func(ctx *Context) []byte

	// ReadCache allows for cache engine replacement. By default, gramework uses github.com/VictoriaMetrics/fastcache.
	// ReadCache returns the value and boolean if the value was found and still valid.
	ReadCache func(ctx *Context, key []byte) ([]byte, bool)
	// StoreCache allows for cache engine replacement. By default, gramework uses github.com/VictoriaMetrics/fastcache.
	StoreCache func(ctx *Context, key, value []byte, ttl time.Duration)

	// CacheableHeaders is a list of headers that gramework can cache.
	// Note, that if X-ABC is present both in cacheable and noncacheable header lists,
	// it will not be cached.
	CacheableHeaders []string // slice of canonical header names
	// NonCacheableHeaders is a list of headers that gramework can not cache.
	// Note, that if X-ABC is present both in cacheable and noncacheable header lists,
	// it will not be cached.
	NonCacheableHeaders []string
}

func (opts *CacheOptions) validate() error {
	if opts.TTL <= 0 {
		return errors.New("TTL must be grater than 0")
	}
	if opts.CacheKey == nil {
		opts.CacheKey = defaultCacheOpts.CacheKey
	}
	if opts.Cacheable == nil {
		opts.Cacheable = defaultCacheOpts.Cacheable
	}

	return nil
}

var defaultCacheOpts = NewCacheOptions()

// NewCacheOptions returns a cache options with default settings.
func NewCacheOptions() *CacheOptions {
	return &CacheOptions{
		TTL: 30 * time.Second,
		Cacheable: func(ctx *Context) bool {
			if len(ctx.Request.Header.Peek("Authentication")) > 0 {
				return false
			}

			if len(ctx.Cookies.Storage) > 0 {
				return false
			}

			return true
		},
		CacheKey: func(ctx *Context) []byte {
			return ctx.Path()
		},
	}
}

// CacheFor is a shortcut to set ttl easily. See app.Cache() for docs.
func (app *App) CacheFor(handler interface{}, ttl time.Duration) func(ctx *Context) {
	opts := app.getCacheOpts()

	opts.TTL = ttl
	return app.Cache(handler, opts)
}

// Cache wrapper will cache given handler using provided options. If options parameter omitted,
// this function will use default options.
//
// NOTE: Please, your CacheOptions' TTL must be more than 0.
func (app *App) Cache(handler interface{}, options ...*CacheOptions) func(ctx *Context) {
	opts := app.getCacheOpts(options...)

	if err := opts.validate(); err != nil {
		app.Logger.WithError(err).Fatal("could not initialize cache middleware: check options")
	}

	wrappedHandler := app.defaultRouter.determineHandler(handler)

	if opts.ReadCache == nil || opts.StoreCache == nil {
		cache := fastcache.New(1)
		opts.ReadCache = readFastCache(cache)
		opts.StoreCache = storeFastCache(cache)
	}

	return func(ctx *Context) {
		if opts.Cacheable(ctx) {
			cacheKey := opts.CacheKey(ctx)
			if value, isValid := opts.ReadCache(ctx, cacheKey); isValid {
				serializedHeaders, isValid := opts.ReadCache(ctx, append(cacheKey, []byte("-headers")...))
				if isValid {
					headers := map[string]string{}
					err := json.Unmarshal(serializedHeaders, &headers)
					if err == nil {
						for name, value := range headers {
							ctx.Response.Header.Set(name, value)
						}
						ctx.Response.SetBody(value)
						return
					}
				}
			}

			wrappedHandler(ctx)

			b := ctx.Response.Body()

			opts.StoreCache(ctx, cacheKey, b, opts.TTL)
			headers, ok := serializeHeaders(ctx, opts)
			if ok {
				opts.StoreCache(ctx, append(cacheKey, []byte("-headers")...), headers, opts.TTL)
			}
			return
		}

		wrappedHandler(ctx)
	}
}

func serializeHeaders(ctx *Context, opts *CacheOptions) ([]byte, bool) {
	headers := map[string]string{
		"Content-Type":   string(ctx.Response.Header.Peek("Content-Type")),
		"Content-Length": string(ctx.Response.Header.Peek("Content-Length")),
	}
	for _, header := range opts.CacheableHeaders {
		headers[header] = string(ctx.Response.Header.Peek(header))
	}
	for _, header := range opts.NonCacheableHeaders {
		delete(headers, header)
	}
	serialized, err := json.Marshal(headers)
	return serialized, err == nil
}

func readFastCache(cache *fastcache.Cache) func(_ *Context, key []byte) (value []byte, isValid bool) {
	return func(_ *Context, key []byte) ([]byte, bool) {
		return cache.GetWithTimeout(nil, key)
	}
}

func storeFastCache(cache *fastcache.Cache) func(_ *Context, key, value []byte, ttl time.Duration) {
	return func(_ *Context, key, value []byte, ttl time.Duration) {
		cache.SetWithTimeout(key, value, ttl)
	}
}

func (app *App) getCacheOpts(options ...*CacheOptions) *CacheOptions {
	opts := defaultCacheOpts
	switch {
	case len(options) > 1:
		app.Logger.Warn("got more than one set of cache options: using the first one.")
		fallthrough
	case len(options) == 1:
		if options[0] != nil {
			opts = options[0]
		}
	case app.DefaultCacheOptions != nil:
		opts = app.DefaultCacheOptions
	}

	return opts
}
