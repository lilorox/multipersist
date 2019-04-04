package main

import (
	"testing"
)

var nbLength = 100

func BenchmarkIncrement(b *testing.B) {
	n := NewNumber(100)

	for i := 0; i < b.N; i++ {
		if !n.Increment() {
			n = NewNumber(100)
		}
	}
}

func benchPersistenceWithCache(size int, b *testing.B) {
	initCache(size)
	n := NewNumber(100)

	for i := 0; i < b.N; i++ {
		n.Persistence()
	}
}

func BenchmarkPersistenceCache5(b *testing.B) {
	benchPersistenceWithCache(5, b)
}

func BenchmarkPersistenceCache10(b *testing.B) {
	benchPersistenceWithCache(10, b)
}

func BenchmarkPersistenceCache15(b *testing.B) {
	benchPersistenceWithCache(15, b)
}
