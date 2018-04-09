package test

import (
	"testing"
)

func errCheck(t *testing.T, err error) {
	if err != nil {
		t.Error(err.Error())
	}
}
