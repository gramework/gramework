package healthchecks

import (
	"errors"
	"runtime"
	"strings"

	sigar "github.com/cloudfoundry/gosigar"
	"github.com/gramework/gramework"
	"github.com/gramework/gramework/internal/gfmt"
)

type hc struct {
	CPUClock string  `json:"cpu_clock"`
	RAM      ramJSON `json:"ram_usage"`
	Swap     ramJSON `json:"swap_usage"`

	LA         laJSON `json:"load_average"`
	LoadStatus string `json:"load_alert_status"`

	Uptime string `json:"uptime"`

	Custom map[string]interface{} `json:"custom_metrics,omitempty"`
}

type laJSON struct {
	One     float64 `json:"one"`
	Five    float64 `json:"five"`
	Fifteen float64 `json:"fifteen"`
}

type ramJSON struct {
	Used  string `json:"used"`
	Total string `json:"total"`
}

type sigarWrapper struct {
	sigar.ConcreteSigar
}

func (s sigarWrapper) swap() ramJSON {
	swap, err := s.GetSwap()
	if err != nil {
		swap.Get()
	}
	return ramJSON{
		Used:  gfmt.Si(swap.Used),
		Total: gfmt.Si(swap.Total),
	}
}
func (s sigarWrapper) ram() ramJSON {
	mem, err := s.GetMem()
	if err != nil {
		mem.Get()
	}
	return ramJSON{
		Used:  gfmt.Si(mem.Used),
		Total: gfmt.Si(mem.Total),
	}
}

func doReg(r interface{}, collectors []func() (statKey string, stats interface{}), registerPing, registerHC bool) error {
	app, isApp := r.(*gramework.App)
	sr, isSr := r.(*gramework.SubRouter)
	if !isApp && !isSr {
		return errors.New("unsupported handler type")
	}

	if isApp {
		sr = app.Sub("")
	}

	if registerPing {
		sr.GET("/ping", "pong")
	}
	if registerHC {
		sr.GET("/healthcheck", ServeHealthcheck(collectors...))
	}

	return nil
}

func check(collectors ...func() (statKey string, stats interface{})) interface{} {
	s := sigarWrapper{sigar.ConcreteSigar{}}
	currentCheck := &hc{
		CPUClock:   gfmt.Si(uint64(gramework.TicksPerSecond())),
		RAM:        s.ram(),
		Swap:       s.swap(),
		LoadStatus: "<unknown>",
	}
	la, err := s.GetLoadAverage()
	if err != nil {
		err = la.Get() // retry
	}

	if err == nil {
		maxLA := float64(runtime.NumCPU() + 2)
		currentCheck.LA = laJSON(la)
		diffOne := maxLA - la.One
		diffFive := maxLA - la.Five
		diffFifteen := maxLA - la.Fifteen

		alertTrigger := float64(-3)
		warnTrigger := float64(0)

		if diffOne < alertTrigger || diffFive < alertTrigger || diffFifteen < alertTrigger {
			currentCheck.LoadStatus = "alert"
		} else if diffOne < warnTrigger || diffFive < warnTrigger || diffFifteen < warnTrigger {
			currentCheck.LoadStatus = "warn"
		} else {
			currentCheck.LoadStatus = "ok"
		}
	}

	uptime := sigar.Uptime{}
	err = uptime.Get()
	if err != nil {
		err = uptime.Get() // retry
	}
	if err == nil {
		currentCheck.Uptime = strings.TrimSpace(uptime.Format())
	}

	if len(collectors) > 0 {
		currentCheck.Custom = make(map[string]interface{})
	}
	for _, cb := range collectors {
		if cb != nil {
			key, stats := cb()
			currentCheck.Custom[key] = stats
		}
	}
	return currentCheck
}
