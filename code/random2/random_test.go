package main

import (
	crnd "crypto/rand"
	oldrand "math/rand"
	"math/rand/v2"
	"testing"
)

func makeSeedChaCha8() [32]byte {
	var seed [32]byte
	_, _ = crnd.Read(seed[:]) //エラー処理をサボってます 🙇
	return seed
}

func makeSeedPPCG() (uint64, uint64) {
	return rand.Uint64(), rand.Uint64()
}

func makeSeedLaggedFibonacci() int64 {
	return rand.Int64()
}

var seedChaCha8 = makeSeedChaCha8()
var seedPCG1, seedPCG2 = makeSeedPPCG()
var seedLaggedFibonacci = makeSeedLaggedFibonacci()
var count = 1000000

func BenchmarkRandomChaCha8(b *testing.B) {
	rnd := rand.New(rand.NewChaCha8(seedChaCha8))
	for i := 0; i < b.N; i++ {
		for range count {
			_ = rnd.IntN(1000)
		}
	}
}

func BenchmarkRandomChaCha8runtime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for range count {
			_ = rand.IntN(1000)
		}
	}
}

func BenchmarkRandomPCG(b *testing.B) {
	rnd := rand.New(rand.NewPCG(seedPCG1, seedPCG2))
	for i := 0; i < b.N; i++ {
		for range count {
			_ = rnd.IntN(1000)
		}
	}
}

func BenchmarkRandomLaggedFibonacci(b *testing.B) {
	rnd := rand.New(oldrand.NewSource(seedLaggedFibonacci).(rand.Source))
	for i := 0; i < b.N; i++ {
		for range count {
			_ = rnd.IntN(1000)
		}
	}
}
