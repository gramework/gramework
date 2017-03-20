// Copyright 2015-2017 Matt Holt and caddy contributors
// Copyright 2017 Kirill Danshin and gramework contributors

package gramework

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/valyala/fasthttp"
)

const challengeBasePath = "/.well-known/acme-challenge"

// HTTPChallengeHandler proxies challenge requests to ACME client if the
// request path starts with challengeBasePath. It returns true if it
// handled the request and no more needs to be done; it returns false
// if this call was a no-op and the request still needs handling.
func HTTPChallengeHandler(ctx *Context, listenHost, altPort string) bool {
	if !strings.HasPrefix(string(ctx.URI().Path()), challengeBasePath) {
		return false
	}
	if DisableHTTPChallenge {
		return false
	}
	if !namesObtaining.Has(string(ctx.Host())) {
		return false
	}

	scheme := "http"
	if ctx.IsTLS() {
		scheme = "https"
	}

	if listenHost == "" {
		listenHost = "localhost"
	}

	upstream, err := url.Parse(fmt.Sprintf("%s://%s:%s", scheme, listenHost, altPort))
	if err != nil {
		ctx.Err500()
		log.Printf("[ERROR] ACME proxy handler: %v", err)
		return true
	}

	proxy := httputil.NewSingleHostReverseProxy(upstream)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	ReverseProxyHandler(ctx, &fasthttp.HostClient{
		Addr: upstream.String(),
	}, upstream)

	return true
}

// ReverseProxyHandler proxies a request to the upstream via given proxyClient
func ReverseProxyHandler(ctx *Context, proxyClient *fasthttp.HostClient, upstream *url.URL) {
	req := &ctx.Request
	resp := &ctx.Response
	prepareRequest(req, upstream)
	if err := proxyClient.Do(req, resp); err != nil {
		ctx.Logger.Errorf("error when proxying the request: %s", err)
	}
	postprocessResponse(resp)
}

func prepareRequest(req *fasthttp.Request, upstream *url.URL) {
	// do not proxy "Connection" header.
	req.Header.Del("Connection")
	req.SetHost(upstream.Host)
	// strip other unneeded headers.

	// alter other request params before sending them to upstream host
}

func postprocessResponse(resp *fasthttp.Response) {
	// do not proxy "Connection" header
	resp.Header.Del("Connection")

	// strip other unneeded headers

	// alter other response data if needed
}
