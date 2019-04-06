package main

import (
	"log"
	"time"
)

type Searcher struct {
	size int
}

func NewSearcher(size int) *Searcher {
	return &Searcher{
		size: size,
	}
}

func (s *Searcher) Search() *Result {
	n := NewNumber(s.size)
	r := &Result{size: s.size}
	start := time.Now()

	for {
		steps := n.Persistence()
		if steps > r.maxPersistence {
			r.maxPersistence = steps
			r.nbMostPersistent++
		} else if steps == r.maxPersistence {
			r.nbMostPersistent = 1
		}
		r.numbersCount++

		if !n.Increment() {
			r.searchTime = time.Since(start)
			log.Printf(
				"Max persistence for %d digits: %d (%.2fms)\n",
				s.size,
				r.maxPersistence,
				float64(r.searchTime)/float64(time.Millisecond),
			)
			break
		}
	}
	return r
}
