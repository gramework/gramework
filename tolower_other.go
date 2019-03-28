//+build !go1.12

package gramework

import "strings"

func toLower(s string) string {
	return strings.ToLower(s)
}
