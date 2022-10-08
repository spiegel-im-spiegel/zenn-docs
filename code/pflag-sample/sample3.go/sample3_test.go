package sample3

import (
	"flag"
	"testing"
)

var foo = flag.Bool("foo", false, "option foo")

func TestMain(m *testing.M) {
	flag.Parse()
	m.Run()
}

func TestFlag(t *testing.T) {
	if !*foo {
		t.Errorf("option foo = %v, want %v.", *foo, true)
	}
}
