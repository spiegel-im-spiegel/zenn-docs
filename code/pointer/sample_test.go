package pointer

import (
	"fmt"
	"testing"
)

type S struct {
	a, b, c int64
	d, e, f string
	g, h, i float64
}

func (s S) Stack(s1 S)    {}
func (s *S) Heap(s1 *S)   {}
func (s S) ValueA() int64 { return s.a }

func byCopy() S {
	return S{
		a: 1, b: 1, c: 1,
		d: "foo", e: "foo", f: "foo",
		g: 1.0, h: 1.0, i: 1.0,
	}
}

func byPointer() *S {
	return &S{
		a: 1, b: 1, c: 1,
		d: "foo", e: "foo", f: "foo",
		g: 1.0, h: 1.0, i: 1.0,
	}
}

func BenchmarkMemoryStack(b *testing.B) {
	var s S
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = byCopy()
	}
	b.StopTimer()
	_ = fmt.Sprintf("%v", s.a)
}

func BenchmarkMemoryHeap(b *testing.B) {
	var s *S
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = byPointer()
	}
	b.StopTimer()
	_ = fmt.Sprintf("%v", s.a)
}

func BenchmarkMemoryStack2(b *testing.B) {
	s, s1 := byCopy(), byCopy()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < 100000; i++ {
			s.Stack(s1)
		}
	}
	b.StopTimer()
	_ = fmt.Sprintf("%v", s.a)
}

func BenchmarkMemoryHeap2(b *testing.B) {
	s, s1 := byPointer(), byPointer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < 100000; i++ {
			s.Heap(s1)
		}
	}
	b.StopTimer()
	_ = fmt.Sprintf("%v", s.a)
}

type IS interface {
	ValueA() int64
}

func byInterface() IS {
	return S{
		a: 1, b: 1, c: 1,
		e: "foo", f: "foo",
		g: 1.0, h: 1.0, i: 1.0,
	}
}

func byInterfaceP() IS {
	return &S{
		a: 1, b: 1, c: 1,
		e: "foo", f: "foo",
		g: 1.0, h: 1.0, i: 1.0,
	}
}

func BenchmarkMemoryInterface(b *testing.B) {
	var s IS

	for i := 0; i < b.N; i++ {
		s = byInterface()
	}
	b.StopTimer()
	_ = fmt.Sprintf("%v", s.ValueA())
}

func BenchmarkMemoryInterfaceP(b *testing.B) {
	var s IS

	for i := 0; i < b.N; i++ {
		s = byInterfaceP()
	}
	b.StopTimer()
	_ = fmt.Sprintf("%v", s.ValueA())
}
