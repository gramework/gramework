package pprof

import (
	"net/http/pprof"

	"github.com/gramework/gramework"
)

var (
	cmdline = gramework.NewGrameHandlerFunc(pprof.Cmdline)
	profile = gramework.NewGrameHandlerFunc(pprof.Profile)
	symbol  = gramework.NewGrameHandlerFunc(pprof.Symbol)
	trace   = gramework.NewGrameHandlerFunc(pprof.Trace)
	index   = gramework.NewGrameHandlerFunc(pprof.Index)
)

// Handler serves server runtime profiling data in the format expected by the pprof visualization tool.
//
// See https://golang.org/pkg/net/http/pprof/ for details.
func Handler(ctx *gramework.Context) {
	ctx.HTML()
	switch ctx.RouteArg("type") {
	case "cmdline":
		cmdline(ctx)
	case "symbol":
		symbol(ctx)
	case "profile":
		profile(ctx)
	case "trace":
		trace(ctx)
	default:
		index(ctx)
	}
}
