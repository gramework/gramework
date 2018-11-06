package gfmt

import "fmt"

// Si formats a number in a short si format
func Si(n uint64) string {
	suff := siRaw
	x := float64(n)
	for ; x >= 1024; x = x / 1024 {
		suff++
	}

	return fmt.Sprintf("%.2f%s", x, suff.String())
}

type siSuff uint

const (
	siRaw siSuff = iota
	siKilo
	siMega
	siGiga
	siTera
)

func (s siSuff) String() string {
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
