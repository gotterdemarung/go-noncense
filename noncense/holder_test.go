package noncense

import (
	"runtime"
	"strconv"
	"testing"
)

const count = 10000

func BenchmarkNoncesAdder_Add(b *testing.B) {
	box := NewNoncesAdder(count)
	runtime.GC()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = <-box.Add(strconv.Itoa(i))
	}
}

func BenchmarkNoncesAdder_Add_Goroutine(b *testing.B) {
	box := NewNoncesAdder(count)
	runtime.GC()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		go func(x int) {
			_ = <-box.Add(strconv.Itoa(i))
		}(i)
	}
}

func BenchmarkNoncesAdderNative_Add(b *testing.B) {
	box := NewNoncesAdderNative(count)
	runtime.GC()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		box.AddSync(strconv.Itoa(i))
	}
}

func BenchmarkNoncesHolder_Add(b *testing.B) {
	box := NewNoncesHolder(uint32(count/2), uint32(count))
	runtime.GC()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		box.Add(NewHString(strconv.Itoa(i)))
	}
}

func BenchmarkHolder_Add(b *testing.B) {
	holder, _ := NewHolder(uint(count))
	runtime.GC()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		holder.Add(strconv.Itoa(i))
	}
}

func BenchmarkHolder_Add_Goroutine(b *testing.B) {
	holder, _ := NewHolder(uint(count))
	runtime.GC()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = <-holder.AddAsync(strconv.Itoa(i))
	}
}
