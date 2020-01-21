package gramework

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/apex/log"
	sigar "github.com/cloudfoundry/gosigar"
)

var checkOnce sync.Once

func checks() {
	checkOnce.Do(func() {
		internalLog.WithFields(log.Fields{
			"package": "gramework",
			"version": Version,
		}).Info("Initialization")
		concreteSigar := sigar.ConcreteSigar{}
		internalLog.WithFields(log.Fields{
			"cputicks": siFmt(uint64(TicksPerSecond())),
			"ram": func() string {
				mem, err := concreteSigar.GetMem()
				if err != nil {
					return "<unknown>"
				}
				return fmt.Sprintf("%s used / %s total", siFmt(mem.ActualUsed), siFmt(mem.Total))
			}(),
			"swap": func() string {
				swap, err := concreteSigar.GetSwap()
				if err != nil {
					return "<unknown>"
				}
				return fmt.Sprintf("%s used / %s total", siFmt(swap.Used), siFmt(swap.Total))
			}(),
		}).Info("node info")
		la, err := concreteSigar.GetLoadAverage()
		if err != nil {
			err = la.Get() // retry
		}
		if err == nil {
			maxLA := float64(runtime.NumCPU() + 2)
			loadLog := internalLog.WithFields(log.Fields{
				"one":     fmt.Sprintf("%.3f", la.One),
				"five":    fmt.Sprintf("%.3f", la.Five),
				"fifteen": fmt.Sprintf("%.3f", la.Fifteen),
			})
			if la.One >= maxLA || la.Five >= maxLA || la.Fifteen >= maxLA {
				loadLog.Warn("high load average, performance may be impacted")
			} else {
				loadLog.Info("load average is good")
			}
		}

		uptime := sigar.Uptime{}
		err = uptime.Get()
		if err != nil {
			err = uptime.Get() // retry
		}
		if err == nil {
			internalLog.WithField("uptime", uptime.Format()).Info("node uptime")
		}
	})
}

func siFmt(n uint64) string {
	prefix := siRaw
	x := float64(n)
	for ; x > 1000; x = x / 1024 {
		prefix++
	}

	return fmt.Sprintf("%.2f%s", x, prefix.String())
}

type siPrefix uint

const (
	siRaw siPrefix = iota
	siKilo
	siMega
	siGiga
	siTera
)

func (s siPrefix) String() string {
	switch s {
	case siRaw:
		return ""
	case siKilo:
		return "K"
	case siMega:
		return "M"
	case siGiga:
		return "G"
	case siTera:
		return "T"
	default:
		return "T"
	}
}
