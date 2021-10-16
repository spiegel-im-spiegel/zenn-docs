//go:build run
// +build run

package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	mrand "math/rand"
)

type Source struct{}

// Seed method is dummy function for rand.Source interface.
func (s Source) Seed(seed int64) {}

// Uint64 method generates a random number in the range [0, 1<<64).
func (s Source) Uint64() uint64 {
	b := [8]byte{}
	ct, _ := rand.Read(b[:])
	return binary.BigEndian.Uint64(b[:ct])
}

// Int63 method generates a random number in the range [0, 1<<63).
func (s Source) Int63() int64 {
	return (int64)(s.Uint64() >> 1)
}

func main() {
	fmt.Println(mrand.New(Source{}).Float64())
}
