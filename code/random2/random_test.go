package randoms

import (
	crnd "crypto/rand"
	oldrand "math/rand"
	"math/rand/v2"
	"testing"

	"github.com/goark/mt/mt19937"
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
var seed1, seed2 = makeSeedPPCG()
var seed3 = makeSeedLaggedFibonacci()

func BenchmarkRandomChaCha8(b *testing.B) {
	rnd := rand.New(rand.NewChaCha8(seedChaCha8))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rnd.IntN(1000)
	}
}

func BenchmarkRandomChaCha8runtime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = rand.IntN(1000)
	}
}

func BenchmarkRandomPCG(b *testing.B) {
	rnd := rand.New(rand.NewPCG(seed1, seed2))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rnd.IntN(1000)
	}
}

func BenchmarkRandomLaggedFibonacci(b *testing.B) {
	rnd := rand.New(oldrand.NewSource(seed3).(rand.Source))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rnd.IntN(1000)
	}
}

func BenchmarkRandomMT(b *testing.B) {
	rnd := rand.New(mt19937.New(seed3))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rnd.IntN(1000)
	}
}
