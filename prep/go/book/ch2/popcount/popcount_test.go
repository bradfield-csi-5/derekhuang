package popcount

import "testing"

func BenchmarkPopCount(b *testing.B) {
	for n := 0; n < b.N; n++ {
		PopCount(0b11110101)
	}
}

func BenchmarkPopCountLoop(b *testing.B) {
	for n := 0; n < b.N; n++ {
		PopCountLoop(0b11110101)
	}
}

func BenchmarkPopCountShift(b *testing.B) {
	for n := 0; n < b.N; n++ {
		PopCountShift(0b11110101)
	}
}

func BenchmarkPopCountClearRightBit(b *testing.B) {
	for n := 0; n < b.N; n++ {
		PopCountClearRightBit(0b11110101)
	}
}
