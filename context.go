package gramework

import (
	"errors"
	"fmt"
)

var (
	// ErrArgNotFound used when no route argument is found
	ErrArgNotFound = errors.New("undefined argument")
)

// Writef is a fmt.Fprintf(context, format, a...) shortcut
func (c *Context) Writef(format string, a ...interface{}) {
	fmt.Fprintf(c, format, a...)
}

// Writeln is a fmt.Fprintln(context, format, a...) shortcut
func (c *Context) Writeln(a ...interface{}) {
	fmt.Fprintln(c, a...)
}

// RouteArg returns an argument value as a string or empty string
func (c *Context) RouteArg(argName string) string {
	v, err := c.RouteArgErr(argName)
	if err != nil {
		return emptyString
	}
	return v
}

// RouteArgErr returns an argument value as a string or empty string
// and ErrArgNotFound if argument was not found
func (c *Context) RouteArgErr(argName string) (string, error) {
	i := c.UserValue(argName)
	if i == nil {
		return emptyString, ErrArgNotFound
	}
	switch value := i.(type) {
	case string:
		return value, nil
	default:
		return fmt.Sprintf(fmtV, i), nil
	}
}

// HTML sets HTML content type
func (c *Context) HTML() *Context {
	c.SetContentType(htmlCT)
	return c
}
