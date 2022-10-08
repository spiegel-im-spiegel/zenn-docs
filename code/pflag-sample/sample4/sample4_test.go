package sample4

import (
	"flag"
	"testing"

	"github.com/spf13/pflag"
)

var foo = pflag.BoolP("foo", "f", false, "option foo")

func TestMain(m *testing.M) {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	flag.Parse()
	m.Run()
}

func TestFlag(t *testing.T) {
	if !*foo {
		t.Errorf("option foo = %v, want %v.", *foo, true)
	}
}
